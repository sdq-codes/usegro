package repositories

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/sdq-codes/usegro-api/internal/apps/form/models"
	dynamodb2 "github.com/sdq-codes/usegro-api/internal/helper/dynamodb"
	"sort"
	"strings"
	"time"
)

type FormRepositoryInterface interface {
	CreateForm(
		ctx context.Context,
		formDynamo *dynamodb.Client,
		version models.FormVersion,
		form models.Form,
	) error
	CreateFormVersion(
		ctx context.Context,
		formDynamo *dynamodb.Client,
		version models.FormVersion,
	) error
	CreateFormVersionField(
		ctx context.Context,
		formDynamo *dynamodb.Client,
		formVersion models.FormVersion,
		field models.FormVersionField,
	) error
	UpdateFormVersionFieldOrder(
		ctx context.Context,
		formDynamo *dynamodb.Client,
		field models.FormVersionField,
	) error
	FetchForm(
		ctx context.Context,
		formDynamo *dynamodb.Client,
		formId string,
	) (*models.CompleteForm, error)

	FetchDraftForm(
		ctx context.Context,
		formDynamo *dynamodb.Client,
		formId string,
	) (*models.CompleteForm, error)

	FetchFormVersion(
		ctx context.Context,
		formDynamo *dynamodb.Client,
		formId string,
		formVersionId string,
	) (*models.CompleteForm, error)
	PublishFormVersion(
		ctx context.Context,
		formDynamo *dynamodb.Client,
		formID string,
		versionID string,
	) error
	DeleteFormVersionField(
		ctx context.Context,
		client *dynamodb.Client,
		formVersion models.FormVersion,
		fieldID string,
	) error

	UpdateFormVersionField(
		ctx context.Context,
		formDynamo *dynamodb.Client,
		field models.FormVersionField,
		updates map[string]interface{},
	) error
}

type FormRepository struct {
	tableName string
}

func NewFormRepository(tableName string) FormRepositoryInterface {
	return &FormRepository{tableName: tableName}
}

func (f *FormRepository) CreateForm(
	ctx context.Context,
	formDynamo *dynamodb.Client,
	version models.FormVersion,
	form models.Form,
) error {
	formItem, err := attributevalue.MarshalMap(form)
	if err != nil {
		return fmt.Errorf("failed to marshal Form: %w", err)
	}

	version.CreatedAt = time.Now()
	version.UpdatedAt = version.CreatedAt
	version.FormVersionStatus = "draft"

	// Marshal version into DynamoDB item
	versionItem, err := attributevalue.MarshalMap(version)
	if err != nil {
		return fmt.Errorf("failed to marshal FormVersion: %w", err)
	}

	// Prepare transact write items
	transactItems := []types.TransactWriteItem{
		{
			Put: &types.Put{
				TableName:           awsString("forms"), // replace with your table name
				Item:                formItem,
				ConditionExpression: awsString("attribute_not_exists(PK) AND attribute_not_exists(SK)"),
			},
		},
		{
			Put: &types.Put{
				TableName:           awsString("forms"),
				Item:                versionItem,
				ConditionExpression: awsString("attribute_not_exists(PK) AND attribute_not_exists(SK)"),
			},
		},
	}

	_, err = formDynamo.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
		TransactItems: transactItems,
	})
	if err != nil {
		return fmt.Errorf("failed to create form version transaction: %w", err)
	}

	return nil
}

func (f *FormRepository) CreateFormVersion(ctx context.Context, formDynamo *dynamodb.Client, version models.FormVersion) error {

	version.CreatedAt = time.Now()
	version.UpdatedAt = version.CreatedAt
	version.FormVersionStatus = "draft"

	// Marshal version into DynamoDB item
	versionItem, err := attributevalue.MarshalMap(version)
	if err != nil {
		return fmt.Errorf("failed to marshal FormVersion: %w", err)
	}

	// Prepare transact write items
	transactItems := []types.TransactWriteItem{
		{
			Put: &types.Put{
				TableName:           awsString("forms"),
				Item:                versionItem,
				ConditionExpression: awsString("attribute_not_exists(PK) AND attribute_not_exists(SK)"),
			},
		},
	}

	_, err = formDynamo.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
		TransactItems: transactItems,
	})
	if err != nil {
		return fmt.Errorf("failed to create form version transaction: %w", err)
	}

	return nil
}

func (f *FormRepository) CreateFormVersionField(
	ctx context.Context,
	formDynamo *dynamodb.Client,
	formVersion models.FormVersion,
	field models.FormVersionField,
) error {
	fieldID := uuid.New().String()
	field.PK = fmt.Sprintf("FORM#%s", dynamodb2.ExtractSK(formVersion.PK))
	field.SK = fmt.Sprintf("VERSION#%sFIELD#%s", dynamodb2.ExtractSK(formVersion.SK), fieldID)
	field.FormVersionID = formVersion.SK
	field.CreatedAt = time.Now().UTC()
	field.UpdatedAt = field.CreatedAt

	// Marshal field to DynamoDB item
	fieldItem, err := attributevalue.MarshalMap(field)
	if err != nil {
		return fmt.Errorf("failed to marshal FormVersionField: %w", err)
	}

	// Insert with condition to prevent overwrite
	_, err = formDynamo.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String("forms"),
		Item:                fieldItem,
		ConditionExpression: aws.String("attribute_not_exists(PK) AND attribute_not_exists(SK)"),
	})
	if err != nil {
		return fmt.Errorf("failed to insert form version field: %w", err)
	}

	return nil
}

func (f *FormRepository) FetchForm(
	ctx context.Context,
	formDynamo *dynamodb.Client,
	formId string,
) (*models.CompleteForm, error) {
	pk := fmt.Sprintf("FORM#%s", formId)

	out, err := formDynamo.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String("forms"),
		KeyConditionExpression: aws.String("PK = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: pk},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query form %s: %w", formId, err)
	}
	if len(out.Items) == 0 {
		return nil, fmt.Errorf("form %s not found", formId)
	}

	fullForm := &models.CompleteForm{}
	var published, drafts, archived []models.FormVersion
	var fields []models.FormVersionField
	var logics []models.FieldLogic

	for _, item := range out.Items {
		var sk string
		_ = attributevalue.Unmarshal(item["SK"], &sk)

		switch {
		// VERSION item
		case strings.HasPrefix(sk, "VERSION#") && !strings.Contains(sk, "FIELD#"):
			var v models.FormVersion
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

		// FIELD item
		case strings.Contains(sk, "VERSION#") && strings.Contains(sk, "FIELD#"):
			var fld models.FormVersionField
			if err := attributevalue.UnmarshalMap(item, &fld); err == nil {
				fields = append(fields, fld)
			}

		// LOGIC item
		case strings.Contains(sk, "LOGIC#"):
			var lg models.FieldLogic
			if err := attributevalue.UnmarshalMap(item, &lg); err == nil {
				logics = append(logics, lg)
			}
		}
	}

	// pick latest version
	var candidates []models.FormVersion
	switch {
	case len(published) > 0:
		candidates = published
	case len(drafts) > 0:
		candidates = drafts
	case len(archived) > 0:
		candidates = archived
	default:
		return nil, fmt.Errorf("no version found for form %s", formId)
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].PublishedAt > candidates[j].PublishedAt
	})
	latestVersion := candidates[0]
	fullForm.Version = latestVersion

	// attach fields + logic
	fieldMap := make(map[string]*models.FormVersionField)
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

	return fullForm, nil
}

func (f *FormRepository) FetchDraftForm(
	ctx context.Context,
	formDynamo *dynamodb.Client,
	formId string,
) (*models.CompleteForm, error) {
	pk := fmt.Sprintf("FORM#%s", formId)

	// 1. Query all items for this Form
	out, err := formDynamo.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String("forms"),
		KeyConditionExpression: aws.String("PK = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: pk},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query form %s: %w", formId, err)
	}
	if len(out.Items) == 0 {
		return nil, fmt.Errorf("form %s not found", formId)
	}

	fullForm := &models.CompleteForm{}
	var drafts []models.FormVersion
	var fields []models.FormVersionField
	var logics []models.FieldLogic

	// 2. Parse items
	for _, item := range out.Items {
		var sk string
		_ = attributevalue.Unmarshal(item["SK"], &sk)

		switch {
		// VERSION item (draft only)
		case strings.HasPrefix(sk, "VERSION#") && !strings.Contains(sk, "FIELD#"):
			var v models.FormVersion
			if err := attributevalue.UnmarshalMap(item, &v); err == nil {
				if v.FormVersionStatus == "draft" {
					drafts = append(drafts, v)
				}
			}

		// FIELD item
		case strings.Contains(sk, "VERSION#") && strings.Contains(sk, "FIELD#"):
			var fld models.FormVersionField
			if err := attributevalue.UnmarshalMap(item, &fld); err == nil {
				fields = append(fields, fld)
			}

		// LOGIC item
		case strings.Contains(sk, "LOGIC#"):
			var lg models.FieldLogic
			if err := attributevalue.UnmarshalMap(item, &lg); err == nil {
				logics = append(logics, lg)
			}
		}
	}

	// 3. Require at least one draft version
	if len(drafts) == 0 {
		return nil, fmt.Errorf("no draft version found for form %s", formId)
	}

	// 4. Pick latest draft (by UpdatedAt)
	sort.Slice(drafts, func(i, j int) bool {
		return drafts[i].UpdatedAt.After(drafts[j].UpdatedAt)
	})
	latestDraft := drafts[0]
	fullForm.Version = latestDraft

	// 5. Attach fields + logic
	fieldMap := make(map[string]*models.FormVersionField)
	for i, fld := range fields {
		if fld.FormVersionID == latestDraft.SK {
			fieldMap[fld.SK] = &fields[i]
		}
	}

	for _, lg := range logics {
		if field, ok := fieldMap[lg.FormVersionFieldID]; ok {
			field.Logic = append(field.Logic, lg)
		}
	}

	for _, f := range fieldMap {
		fullForm.Fields = append(fullForm.Fields, *f)
	}

	// 6. Order fields by "Order" column
	sort.Slice(fullForm.Fields, func(i, j int) bool {
		return fullForm.Fields[i].Order < fullForm.Fields[j].Order
	})

	return fullForm, nil
}

func (f *FormRepository) FetchFormVersion(
	ctx context.Context,
	formDynamo *dynamodb.Client,
	formId string,
	formVersionId string,
) (*models.CompleteForm, error) {
	pk := fmt.Sprintf("FORM#%s", formId)

	// Query all items for this Form (version, fields, logic)
	out, err := formDynamo.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String("forms"),
		KeyConditionExpression: aws.String("PK = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: pk},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query form %s: %w", formId, err)
	}

	if len(out.Items) == 0 {
		return nil, fmt.Errorf("form %s not found", formId)
	}

	result := &models.CompleteForm{}
	var fields []models.FormVersionField
	var logics []models.FieldLogic

	// Parse items
	for _, item := range out.Items {
		var sk string
		_ = attributevalue.Unmarshal(item["SK"], &sk)
		switch {
		// VERSION item
		case dynamodb2.ExtractSK(sk) == formVersionId:
			var v models.FormVersion
			if err := attributevalue.UnmarshalMap(item, &v); err == nil {
				result.Version = v
			}

		// FIELD item (must contain VERSION# + FIELD#)
		case strings.Contains(sk, "VERSION#") && strings.Contains(sk, "FIELD#"):
			var fld models.FormVersionField
			if err := attributevalue.UnmarshalMap(item, &fld); err == nil {
				fieldIds := ExtractFieldIDs(fld.SK)
				if fieldIds[0] == formVersionId {
					fields = append(fields, fld)
				}
			}

		// LOGIC item
		case strings.Contains(sk, "LOGIC#"):
			var lg models.FieldLogic
			if err := attributevalue.UnmarshalMap(item, &lg); err == nil {
				logics = append(logics, lg)
			}
		}
	}

	if result.Version.PK == "" {
		return nil, fmt.Errorf("form version %s not found for form %s", formVersionId, formId)
	}

	// Attach logic to fields
	fieldMap := make(map[string]*models.FormVersionField)
	for i, fld := range fields {
		fieldMap[fld.SK] = &fields[i]
	}

	for _, lg := range logics {
		if f, ok := fieldMap[lg.FormVersionFieldID]; ok {
			f.Logic = append(f.Logic, lg)
		}
	}

	for _, f := range fieldMap {
		result.Fields = append(result.Fields, *f)
	}

	return result, nil
}

func (f *FormRepository) PublishFormVersion(
	ctx context.Context,
	formDynamo *dynamodb.Client,
	formID string,
	versionID string,
) error {
	pk := fmt.Sprintf("FORM#%s", formID)
	targetSk := fmt.Sprintf("VERSION#%s", versionID)

	// 1. Query all versions for this form
	out, err := formDynamo.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String("forms"),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :skprefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":       &types.AttributeValueMemberS{Value: pk},
			":skprefix": &types.AttributeValueMemberS{Value: "VERSION#"},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to query form versions for %s: %w", formID, err)
	}

	var transactItems []types.TransactWriteItem

	// 2. Unpublish all published versions
	for _, item := range out.Items {
		var v models.FormVersion
		if err := attributevalue.UnmarshalMap(item, &v); err != nil {
			continue
		}

		if v.FormVersionStatus == "published" {
			update := types.TransactWriteItem{
				Update: &types.Update{
					TableName: aws.String("forms"),
					Key: map[string]types.AttributeValue{
						"PK": &types.AttributeValueMemberS{Value: pk},
						"SK": &types.AttributeValueMemberS{Value: v.SK},
					},
					UpdateExpression: aws.String("SET formVersionStatus = :formVersionStatus"),
					ExpressionAttributeValues: map[string]types.AttributeValue{
						":formVersionStatus": &types.AttributeValueMemberS{Value: "archived"},
					},
				},
			}
			transactItems = append(transactItems, update)
		}
	}

	// 3. Publish target version
	now := time.Now().UTC().Format(time.RFC3339)
	publishUpdate := types.TransactWriteItem{
		Update: &types.Update{
			TableName: aws.String("forms"),
			Key: map[string]types.AttributeValue{
				"PK": &types.AttributeValueMemberS{Value: pk},
				"SK": &types.AttributeValueMemberS{Value: targetSk},
			},
			UpdateExpression: aws.String("SET formVersionStatus = :formVersionStatus, publishedAt = :publishedAt"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":formVersionStatus": &types.AttributeValueMemberS{Value: "published"},
				":publishedAt":       &types.AttributeValueMemberS{Value: now},
			},
		},
	}
	transactItems = append(transactItems, publishUpdate)

	// 4. Execute transaction
	_, err = formDynamo.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
		TransactItems: transactItems,
	})
	if err != nil {
		return fmt.Errorf("failed to publish form version %s: %w", versionID, err)
	}

	return nil
}

func (f *FormRepository) DeleteFormVersionField(
	ctx context.Context,
	client *dynamodb.Client,
	formVersion models.FormVersion,
	fieldID string,
) error {
	if formVersion.FormVersionStatus != "draft" {
		return errors.New("only draft form versions are supported")
	}

	// Correct SK (matches how fields are created)
	sk := fmt.Sprintf("VERSION#%sFIELD#%s", dynamodb2.ExtractSK(formVersion.SK), fieldID)

	// Delete the field
	_, err := client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String("forms"),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: formVersion.PK},
			"SK": &types.AttributeValueMemberS{Value: sk},
		},
		ConditionExpression: aws.String("attribute_exists(PK) AND attribute_exists(SK)"),
	})
	if err != nil {
		return fmt.Errorf("failed to delete form version field: %w", err)
	}

	return nil
}

func awsString(s string) *string {
	return &s
}

func ExtractFieldIDs(sk string) []string {
	// Example: "VERSION#<formVersionId>FIELD#<fieldId>"
	parts := strings.Split(sk, "FIELD#")
	if len(parts) != 2 {
		return nil
	}

	// Extract formVersionId (after "VERSION#")
	formVersionID := strings.TrimPrefix(parts[0], "VERSION#")
	fieldID := parts[1]

	return []string{formVersionID, fieldID}
}

func (f *FormRepository) UpdateFormVersionFieldOrder(
	ctx context.Context,
	formDynamo *dynamodb.Client,
	field models.FormVersionField,
) error {
	field.UpdatedAt = time.Now().UTC()

	update := expression.Set(
		expression.Name("order"), expression.Value(field.Order),
	).Set(
		expression.Name("updatedAt"), expression.Value(field.UpdatedAt),
	)

	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return fmt.Errorf("failed to build update expression: %w", err)
	}

	_, err = formDynamo.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String("forms"),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: field.PK},
			"SK": &types.AttributeValueMemberS{Value: field.SK},
		},
		UpdateExpression:          expr.Update(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		return fmt.Errorf("failed to update field order: %w", err)
	}

	return nil
}

func (f *FormRepository) UpdateFormVersionField(
	ctx context.Context,
	formDynamo *dynamodb.Client,
	field models.FormVersionField,
	updates map[string]interface{},
) error {

	if len(updates) == 0 {
		return nil // nothing to update
	}

	// Always update UpdatedAt
	updates["UpdatedAt"] = time.Now().UTC()

	// Build update expression
	updateBuilder := expression.UpdateBuilder{}
	for k, v := range updates {
		updateBuilder = updateBuilder.Set(expression.Name(k), expression.Value(v))
	}

	expr, err := expression.NewBuilder().WithUpdate(updateBuilder).Build()
	if err != nil {
		return fmt.Errorf("failed to build update expression: %w", err)
	}

	_, err = formDynamo.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String("forms"),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: field.PK},
			"SK": &types.AttributeValueMemberS{Value: field.SK},
		},
		UpdateExpression:          expr.Update(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		return fmt.Errorf("failed to update form field %s: %w", field.SK, err)
	}

	return nil
}
