package repositories

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/redis/go-redis/v9"
	"github.com/usegro/services/crm/internal/apps/crm/models"
	formsModel "github.com/usegro/services/crm/internal/apps/form/models"
	"gorm.io/gorm"
	"sort"
	"strings"
	"time"
)

type CRMCustomerRepositoryInterface interface {
	FetchCreateCustomerForm(ctx context.Context, tx *gorm.DB, crmId string) (*[]models.CrmUserOrganization, error)
	FetchCrmCustomers(
		ctx context.Context,
		dynamo *dynamodb.Client,
		crmId string,
		formId string,
	) ([]formsModel.FormSubmission, error)
	ArchiveCrmCustomer(
		ctx context.Context,
		dynamo *dynamodb.Client,
		submissionID,
		crmId,
		formId string,
	) error
	GetCrmCustomer(
		ctx context.Context,
		dynamo *dynamodb.Client,
		submissionID,
		formId string,
		crmId string,
	) (*formsModel.FormSubmission, error)
}

type CRMCustomerRepository struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewCRMCustomerRepository(db *gorm.DB, rdb *redis.Client) *CRMCustomerRepository {
	return &CRMCustomerRepository{
		db:  db,
		rdb: rdb,
	}
}

func (c *CRMCustomerRepository) FetchCrmCustomers(
	ctx context.Context,
	dynamo *dynamodb.Client,
	crmId string,
) ([]formsModel.FormSubmission, error) {
	out, err := dynamo.Scan(ctx, &dynamodb.ScanInput{
		TableName:        aws.String("form_submissions"),
		FilterExpression: aws.String("crmID = :crmId AND #t = :type AND attribute_not_exists(archivedAt)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":crmId": &types.AttributeValueMemberS{Value: crmId},
			":type":  &types.AttributeValueMemberS{Value: "customer"},
		},
		ExpressionAttributeNames: map[string]string{
			"#t": "type",
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to scan crm_customers for crmID %s: %w", crmId, err)
	}

	if len(out.Items) == 0 {
		return nil, fmt.Errorf("no customers found for crmID %s", crmId)
	}

	var customers []formsModel.FormSubmission
	for _, item := range out.Items {
		var customer formsModel.FormSubmission
		if err := attributevalue.UnmarshalMap(item, &customer); err != nil {
			continue
		}
		customers = append(customers, customer)
	}

	return customers, nil
}

func (c *CRMCustomerRepository) GetCrmCustomer(
	ctx context.Context,
	dynamo *dynamodb.Client,
	submissionID,
	formId string,
	crmId string,
) (*formsModel.FormSubmission, error) {
	key := map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("FORM#%s", formId)},
		"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("SUBMISSION#%s", submissionID)},
	}

	out, err := dynamo.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String("form_submissions"),
		Key:       key,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to fetch CRM customer %s: %w", submissionID, err)
	}

	if out.Item == nil || len(out.Item) == 0 {
		return nil, fmt.Errorf("no CRM customer found for submission %s", submissionID)
	}

	var customer formsModel.FormSubmission
	if err := attributevalue.UnmarshalMap(out.Item, &customer); err != nil {
		return nil, fmt.Errorf("failed to unmarshal CRM customer: %w", err)
	}

	// ✅ Manual filter logic
	if customer.CrmID != crmId {
		return nil, fmt.Errorf("crmID mismatch for submission %s", submissionID)
	}

	if customer.ArchivedAt != nil {
		return nil, fmt.Errorf("customer %s is archived", submissionID)
	}

	return &customer, nil

}

func (c *CRMCustomerRepository) ArchiveCrmCustomer(
	ctx context.Context,
	dynamo *dynamodb.Client,
	submissionID,
	formId,
	crmId string,
) error {
	now := time.Now().UTC().Format(time.RFC3339)

	key := map[string]types.AttributeValue{
		"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("SUBMISSION#%s", submissionID)},
		"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("FORM#%s", formId)},
	}

	_, err := dynamo.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:           aws.String("form_submissions"),
		Key:                 key,
		UpdateExpression:    aws.String("SET archivedAt = :archivedAt"),
		ConditionExpression: aws.String("crmID = :crmId AND #t = :type"),
		ExpressionAttributeNames: map[string]string{
			"#t": "type",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":archivedAt": &types.AttributeValueMemberS{Value: now},
			":crmId":      &types.AttributeValueMemberS{Value: crmId},
			":type":       &types.AttributeValueMemberS{Value: "customer"},
		},
		ReturnValues: types.ReturnValueUpdatedNew,
	})

	if err != nil {
		return fmt.Errorf("failed to archive submission %s: %w", submissionID, err)
	}

	return nil
}

func (c *CRMCustomerRepository) FetchPublishedCreateCustomerForm(
	ctx context.Context,
	formDynamo *dynamodb.Client,
	crmId string,
) (*formsModel.CompleteForm, error) {
	form, err := c.fetchPublishedCreateCustomerFormByCRM(ctx, formDynamo, crmId)
	if err == nil {
		return form, nil
	}

	// Fallback: try global
	form, err = c.fetchPublishedCreateCustomerFormByCRM(ctx, formDynamo, "global")

	if err != nil {
		return nil, fmt.Errorf("no published create_customer form found for crm %s or global: %w", crmId, err)
	}
	return form, nil
}

func (c *CRMCustomerRepository) fetchPublishedCreateCustomerFormByCRM(
	ctx context.Context,
	formDynamo *dynamodb.Client,
	crmId string,
) (*formsModel.CompleteForm, error) {

	outForm, err := formDynamo.Scan(ctx, &dynamodb.ScanInput{
		TableName:        aws.String("forms"),
		FilterExpression: aws.String("crmID = :crmId AND #t = :type"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":crmId": &types.AttributeValueMemberS{Value: crmId},
			":type":  &types.AttributeValueMemberS{Value: "create_customer"},
		},
		ExpressionAttributeNames: map[string]string{
			"#t": "type",
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to query forms for crmID %s: %w", crmId, err)
	}
	if len(outForm.Items) == 0 {
		return nil, fmt.Errorf("no forms found for crmID %s", crmId)
	}

	fullForm := &formsModel.FullForm{}
	var form formsModel.Form
	var formId string

	if err := attributevalue.UnmarshalMap(outForm.Items[0], &form); err == nil {
		formId = form.PK // capture the formId (or use form.ID if that’s your field)
		fullForm.Form = form
	}

	out, err := formDynamo.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String("forms"),
		KeyConditionExpression: aws.String("PK = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: formId},
		},
	})

	var published, drafts, archived []formsModel.FormVersion
	var fields []formsModel.FormVersionField
	var logics []formsModel.FieldLogic

	// loop through results
	for _, item := range out.Items {
		var sk string
		_ = attributevalue.Unmarshal(item["SK"], &sk)

		switch {
		case strings.HasPrefix(sk, "VERSION#") && !strings.Contains(sk, "FIELD#"):
			var v formsModel.FormVersion
			if err := attributevalue.UnmarshalMap(item, &v); err == nil {
				switch v.FormVersionStatus {
				case "published":
					published = append(published, v)
				case "draft":
					drafts = append(drafts, v)
				case "archived":
					archived = append(archived, v)
				}
			}
		case strings.Contains(sk, "VERSION#") && strings.Contains(sk, "FIELD#"):
			var fld formsModel.FormVersionField
			if err := attributevalue.UnmarshalMap(item, &fld); err == nil {
				fields = append(fields, fld)
			}

		// LOGIC item
		case strings.Contains(sk, "LOGIC#"):
			var lg formsModel.FieldLogic
			if err := attributevalue.UnmarshalMap(item, &lg); err == nil {
				logics = append(logics, lg)
			}
		}
	}

	if formId == "" {
		return nil, fmt.Errorf("no FORM record found for crmID %s", crmId)
	}

	// pick latest version
	var candidates []formsModel.FormVersion
	switch {
	case len(published) > 0:
		candidates = published
	case len(drafts) > 0:
		candidates = drafts
	case len(archived) > 0:
		candidates = archived
	default:
		return nil, fmt.Errorf("no versions found for form %s", formId)
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].PublishedAt > candidates[j].PublishedAt
	})
	latestVersion := candidates[0]
	fullForm.Version = latestVersion

	fieldMap := make(map[string]*formsModel.FormVersionField)
	for i, fld := range fields {
		if fld.FormVersionID == latestVersion.SK {
			fieldMap[fld.SK] = &fields[i]
		}
	}
	for _, lg := range logics {
		if f, ok := fieldMap[lg.FormVersionFieldID]; ok {
			f.Logic = append(f.Logic, lg)
		}
	}

	// collect + sort by Order column
	for _, f := range fieldMap {
		fullForm.Fields = append(fullForm.Fields, *f)
	}
	sort.Slice(fullForm.Fields, func(i, j int) bool {
		return fullForm.Fields[i].Order < fullForm.Fields[j].Order
	})

	return &formsModel.CompleteForm{
		Version: fullForm.Version,
		Fields:  fullForm.Fields,
	}, nil
}
