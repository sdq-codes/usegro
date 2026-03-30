package repositories

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/sdq-codes/usegro-api/internal/apps/form/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FormRepositoryInterface interface {
	CreateForm(ctx context.Context, version models.FormVersion, form models.Form) error
	CreateFormVersion(ctx context.Context, version models.FormVersion) error
	CreateFormVersionField(ctx context.Context, formVersionID string, field models.FormVersionField) error
	UpdateFormVersionFieldOrder(ctx context.Context, field models.FormVersionField) error
	FetchForm(ctx context.Context, formId string) (*models.CompleteForm, error)
	FetchDraftForm(ctx context.Context, formId string) (*models.CompleteForm, error)
	FetchFormVersion(ctx context.Context, formId string, formVersionId string) (*models.CompleteForm, error)
	PublishFormVersion(ctx context.Context, formID string, versionID string) error
	DeleteFormVersionField(ctx context.Context, formVersionID string, fieldID string) error
	UpdateFormVersionField(ctx context.Context, fieldID string, updates map[string]interface{}) error
}

type FormRepository struct {
	db *mongo.Database
}

func NewFormRepository(db *mongo.Database) FormRepositoryInterface {
	return &FormRepository{db: db}
}

func (f *FormRepository) forms() *mongo.Collection    { return f.db.Collection("forms") }
func (f *FormRepository) versions() *mongo.Collection { return f.db.Collection("form_versions") }
func (f *FormRepository) fields() *mongo.Collection   { return f.db.Collection("form_fields") }

func (f *FormRepository) CreateForm(ctx context.Context, version models.FormVersion, form models.Form) error {
	version.CreatedAt = time.Now()
	version.UpdatedAt = version.CreatedAt
	version.FormVersionStatus = "draft"

	session, err := f.db.Client().StartSession()
	if err != nil {
		return fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sc mongo.SessionContext) (interface{}, error) {
		if _, err := f.forms().InsertOne(sc, form); err != nil {
			return nil, fmt.Errorf("failed to insert form: %w", err)
		}
		if _, err := f.versions().InsertOne(sc, version); err != nil {
			return nil, fmt.Errorf("failed to insert form version: %w", err)
		}
		return nil, nil
	})
	return err
}

func (f *FormRepository) CreateFormVersion(ctx context.Context, version models.FormVersion) error {
	version.CreatedAt = time.Now()
	version.UpdatedAt = version.CreatedAt
	version.FormVersionStatus = "draft"

	_, err := f.versions().InsertOne(ctx, version)
	if err != nil {
		return fmt.Errorf("failed to create form version: %w", err)
	}
	return nil
}

func (f *FormRepository) CreateFormVersionField(ctx context.Context, formVersionID string, field models.FormVersionField) error {
	field.ID = uuid.New().String()
	field.FormVersionID = formVersionID
	field.CreatedAt = time.Now().UTC()
	field.UpdatedAt = field.CreatedAt

	_, err := f.fields().InsertOne(ctx, field)
	if err != nil {
		return fmt.Errorf("failed to insert form version field: %w", err)
	}
	return nil
}

func (f *FormRepository) fetchVersionsForForm(ctx context.Context, formId string) ([]models.FormVersion, error) {
	cur, err := f.versions().Find(ctx, bson.M{"formID": formId})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var versions []models.FormVersion
	if err := cur.All(ctx, &versions); err != nil {
		return nil, err
	}
	return versions, nil
}

func (f *FormRepository) fetchFieldsForVersion(ctx context.Context, versionID string) ([]models.FormVersionField, error) {
	opts := options.Find().SetSort(bson.D{{Key: "order", Value: 1}})
	cur, err := f.fields().Find(ctx, bson.M{"formVersionID": versionID}, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var fields []models.FormVersionField
	if err := cur.All(ctx, &fields); err != nil {
		return nil, err
	}
	return fields, nil
}

func pickLatestVersion(versions []models.FormVersion) (*models.FormVersion, error) {
	var published, drafts, archived []models.FormVersion
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

	var candidates []models.FormVersion
	switch {
	case len(published) > 0:
		candidates = published
	case len(drafts) > 0:
		candidates = drafts
	case len(archived) > 0:
		candidates = archived
	default:
		return nil, fmt.Errorf("no version found")
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].PublishedAt > candidates[j].PublishedAt
	})
	return &candidates[0], nil
}

func (f *FormRepository) FetchForm(ctx context.Context, formId string) (*models.CompleteForm, error) {
	versions, err := f.fetchVersionsForForm(ctx, formId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch versions for form %s: %w", formId, err)
	}
	if len(versions) == 0 {
		return nil, fmt.Errorf("form %s not found", formId)
	}

	latest, err := pickLatestVersion(versions)
	if err != nil {
		return nil, fmt.Errorf("no version found for form %s", formId)
	}

	fields, err := f.fetchFieldsForVersion(ctx, latest.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch fields: %w", err)
	}

	return &models.CompleteForm{Version: *latest, Fields: fields}, nil
}

func (f *FormRepository) FetchDraftForm(ctx context.Context, formId string) (*models.CompleteForm, error) {
	versions, err := f.fetchVersionsForForm(ctx, formId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch versions for form %s: %w", formId, err)
	}

	var drafts []models.FormVersion
	for _, v := range versions {
		if v.FormVersionStatus == "draft" {
			drafts = append(drafts, v)
		}
	}
	if len(drafts) == 0 {
		return nil, fmt.Errorf("no draft version found for form %s", formId)
	}

	sort.Slice(drafts, func(i, j int) bool {
		return drafts[i].UpdatedAt.After(drafts[j].UpdatedAt)
	})
	latest := drafts[0]

	fields, err := f.fetchFieldsForVersion(ctx, latest.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch fields: %w", err)
	}

	return &models.CompleteForm{Version: latest, Fields: fields}, nil
}

func (f *FormRepository) FetchFormVersion(ctx context.Context, formId string, formVersionId string) (*models.CompleteForm, error) {
	var version models.FormVersion
	err := f.versions().FindOne(ctx, bson.M{"_id": formVersionId, "formID": formId}).Decode(&version)
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("form version %s not found for form %s", formVersionId, formId)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to fetch form version: %w", err)
	}

	fields, err := f.fetchFieldsForVersion(ctx, version.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch fields: %w", err)
	}

	return &models.CompleteForm{Version: version, Fields: fields}, nil
}

func (f *FormRepository) PublishFormVersion(ctx context.Context, formID string, versionID string) error {
	session, err := f.db.Client().StartSession()
	if err != nil {
		return fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sc mongo.SessionContext) (interface{}, error) {
		// Archive all currently published versions for this form
		_, err := f.versions().UpdateMany(sc,
			bson.M{"formID": formID, "formVersionStatus": "published"},
			bson.M{"$set": bson.M{"formVersionStatus": "archived"}},
		)
		if err != nil {
			return nil, fmt.Errorf("failed to archive published versions: %w", err)
		}

		now := time.Now().UTC().Format(time.RFC3339)
		_, err = f.versions().UpdateOne(sc,
			bson.M{"_id": versionID, "formID": formID},
			bson.M{"$set": bson.M{"formVersionStatus": "published", "publishedAt": now}},
		)
		if err != nil {
			return nil, fmt.Errorf("failed to publish version %s: %w", versionID, err)
		}
		return nil, nil
	})
	return err
}

func (f *FormRepository) DeleteFormVersionField(ctx context.Context, formVersionID string, fieldID string) error {
	// Verify the version is a draft before allowing deletion
	var version models.FormVersion
	err := f.versions().FindOne(ctx, bson.M{"_id": formVersionID}).Decode(&version)
	if err != nil {
		return fmt.Errorf("form version not found: %w", err)
	}
	if version.FormVersionStatus != "draft" {
		return errors.New("only draft form versions are supported")
	}

	res, err := f.fields().DeleteOne(ctx, bson.M{"_id": fieldID, "formVersionID": formVersionID})
	if err != nil {
		return fmt.Errorf("failed to delete form version field: %w", err)
	}
	if res.DeletedCount == 0 {
		return fmt.Errorf("form version field not found")
	}
	return nil
}

func (f *FormRepository) UpdateFormVersionFieldOrder(ctx context.Context, field models.FormVersionField) error {
	field.UpdatedAt = time.Now().UTC()
	_, err := f.fields().UpdateOne(ctx,
		bson.M{"_id": field.ID},
		bson.M{"$set": bson.M{"order": field.Order, "updatedAt": field.UpdatedAt}},
	)
	if err != nil {
		return fmt.Errorf("failed to update field order: %w", err)
	}
	return nil
}

func (f *FormRepository) UpdateFormVersionField(ctx context.Context, fieldID string, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}
	updates["updatedAt"] = time.Now().UTC()
	_, err := f.fields().UpdateOne(ctx,
		bson.M{"_id": fieldID},
		bson.M{"$set": updates},
	)
	if err != nil {
		return fmt.Errorf("failed to update form field %s: %w", fieldID, err)
	}
	return nil
}

// EnsureFormIndexes creates required MongoDB indexes
func EnsureFormIndexes(ctx context.Context, db *mongo.Database) error {
	_, err := db.Collection("forms").Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "crmID", Value: 1}, {Key: "type", Value: 1}}},
	})
	if err != nil {
		return err
	}
	_, err = db.Collection("form_versions").Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "formID", Value: 1}}},
		{Keys: bson.D{{Key: "formID", Value: 1}, {Key: "formVersionStatus", Value: 1}}},
	})
	if err != nil {
		return err
	}
	_, err = db.Collection("form_fields").Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "formVersionID", Value: 1}, {Key: "order", Value: 1}}},
	})
	return err
}
