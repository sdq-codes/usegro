package repositories

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/usegro/services/crm/internal/apps/crm/models"
	formsModel "github.com/usegro/services/crm/internal/apps/form/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
)

type PaginatedCustomers struct {
	Data       []formsModel.FormSubmission `json:"data"`
	Total      int64                       `json:"total"`
	Page       int                         `json:"page"`
	Limit      int                         `json:"limit"`
	TotalPages int                         `json:"total_pages"`
}

type CRMCustomerRepositoryInterface interface {
	FetchCreateCustomerForm(ctx context.Context, tx *gorm.DB, crmId string) (*[]models.CrmUserOrganization, error)
	FetchCrmCustomers(ctx context.Context, crmId string, page, limit int) (*PaginatedCustomers, error)
	ArchiveCrmCustomer(ctx context.Context, submissionID, crmId, formId string) error
	GetCrmCustomer(ctx context.Context, submissionID, formId string, crmId string) (*formsModel.FormSubmission, error)
	FetchPublishedCreateCustomerForm(ctx context.Context, crmId string) (*formsModel.CompleteForm, error)
}

type CRMCustomerRepository struct {
	db      *gorm.DB
	rdb     *redis.Client
	mongoDB *mongo.Database
}

func NewCRMCustomerRepository(db *gorm.DB, rdb *redis.Client, mongoDB *mongo.Database) *CRMCustomerRepository {
	return &CRMCustomerRepository{db: db, rdb: rdb, mongoDB: mongoDB}
}

// FetchCreateCustomerForm uses Postgres (gorm) — unchanged.
func (c *CRMCustomerRepository) FetchCreateCustomerForm(ctx context.Context, tx *gorm.DB, crmId string) (*[]models.CrmUserOrganization, error) {
	return nil, nil
}

func (c *CRMCustomerRepository) FetchCrmCustomers(ctx context.Context, crmId string, page, limit int) (*PaginatedCustomers, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	filter := bson.M{
		"crmID":      crmId,
		"type":       "customer",
		"archivedAt": bson.M{"$exists": false},
	}

	total, err := c.mongoDB.Collection("form_submissions").CountDocuments(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to count customers for crmID %s: %w", crmId, err)
	}

	skip := int64((page - 1) * limit)
	opts := options.Find().
		SetSkip(skip).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cur, err := c.mongoDB.Collection("form_submissions").Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch customers for crmID %s: %w", crmId, err)
	}
	defer cur.Close(ctx)

	customers := make([]formsModel.FormSubmission, 0)
	if err := cur.All(ctx, &customers); err != nil {
		return nil, fmt.Errorf("failed to decode customers: %w", err)
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	if totalPages == 0 {
		totalPages = 1
	}

	return &PaginatedCustomers{
		Data:       customers,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (c *CRMCustomerRepository) GetCrmCustomer(ctx context.Context, submissionID, formId string, crmId string) (*formsModel.FormSubmission, error) {
	var customer formsModel.FormSubmission
	err := c.mongoDB.Collection("form_submissions").FindOne(ctx, bson.M{
		"_id":    submissionID,
		"formID": formId,
	}).Decode(&customer)
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("no CRM customer found for submission %s", submissionID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to fetch CRM customer %s: %w", submissionID, err)
	}
	if customer.CrmID != crmId {
		return nil, fmt.Errorf("crmID mismatch for submission %s", submissionID)
	}
	if customer.ArchivedAt != nil {
		return nil, fmt.Errorf("customer %s is archived", submissionID)
	}
	return &customer, nil
}

func (c *CRMCustomerRepository) ArchiveCrmCustomer(ctx context.Context, submissionID, formId, crmId string) error {
	now := time.Now().UTC()
	res, err := c.mongoDB.Collection("form_submissions").UpdateOne(ctx,
		bson.M{"_id": submissionID, "formID": formId, "crmID": crmId, "type": "customer"},
		bson.M{"$set": bson.M{"archivedAt": now}},
	)
	if err != nil {
		return fmt.Errorf("failed to archive submission %s: %w", submissionID, err)
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("submission %s not found or crmID mismatch", submissionID)
	}
	return nil
}

func (c *CRMCustomerRepository) FetchPublishedCreateCustomerForm(ctx context.Context, crmId string) (*formsModel.CompleteForm, error) {
	form, err := c.fetchPublishedCreateCustomerFormByCRM(ctx, crmId)
	if err == nil {
		return form, nil
	}
	return c.fetchPublishedCreateCustomerFormByCRM(ctx, "global")
}

func (c *CRMCustomerRepository) fetchPublishedCreateCustomerFormByCRM(ctx context.Context, crmId string) (*formsModel.CompleteForm, error) {
	var form formsModel.Form
	err := c.mongoDB.Collection("forms").FindOne(ctx, bson.M{"crmID": crmId, "type": "create_customer"}).Decode(&form)
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("no forms found for crmID %s", crmId)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query forms for crmID %s: %w", crmId, err)
	}

	// Fetch all versions for this form
	cur, err := c.mongoDB.Collection("form_versions").Find(ctx, bson.M{"formID": form.ID})
	if err != nil {
		return nil, fmt.Errorf("failed to query form versions: %w", err)
	}
	defer cur.Close(ctx)

	var versions []formsModel.FormVersion
	if err := cur.All(ctx, &versions); err != nil {
		return nil, fmt.Errorf("failed to decode versions: %w", err)
	}

	// Pick latest version (published > draft > archived)
	var published, drafts, archived []formsModel.FormVersion
	for _, v := range versions {
		switch v.FormVersionStatus {
		case "published":
			published = append(published, v)
		case "draft":
			drafts = append(drafts, v)
		case "archived":
			archived = append(archived, v)
		}
	}
	var candidates []formsModel.FormVersion
	switch {
	case len(published) > 0:
		candidates = published
	case len(drafts) > 0:
		candidates = drafts
	case len(archived) > 0:
		candidates = archived
	default:
		return nil, fmt.Errorf("no versions found for form %s", form.ID)
	}
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].PublishedAt > candidates[j].PublishedAt
	})
	latestVersion := candidates[0]

	// Fetch fields for this version
	fieldOpts := options.Find().SetSort(bson.D{{Key: "order", Value: 1}})
	fieldCur, err := c.mongoDB.Collection("form_fields").Find(ctx, bson.M{"formVersionID": latestVersion.ID}, fieldOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to query fields: %w", err)
	}
	defer fieldCur.Close(ctx)

	var fields []formsModel.FormVersionField
	if err := fieldCur.All(ctx, &fields); err != nil {
		return nil, fmt.Errorf("failed to decode fields: %w", err)
	}

	return &formsModel.CompleteForm{Version: latestVersion, Fields: fields}, nil
}
