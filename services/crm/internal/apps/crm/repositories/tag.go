package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"

	"github.com/usegro/services/crm/internal/apps/crm/models"
)

const TagTableName = "tags" // change this to your actual table name

type TagRepositoryInterface interface {
	CreateTag(ctx context.Context, db *dynamodb.Client, tag models.Tag) (*models.Tag, error)
	FetchTag(ctx context.Context, db *dynamodb.Client, crmID string, tagID string) (*models.Tag, error)
	ListTagsByCRM(ctx context.Context, db *dynamodb.Client, crmID string) ([]models.Tag, error)
	UpdateTagName(ctx context.Context, db *dynamodb.Client, crmID string, tagID string, newName string) error
	UpdateTagStatus(ctx context.Context, db *dynamodb.Client, crmID string, tagID string, status string) error
	ArchiveTag(ctx context.Context, db *dynamodb.Client, crmID string, tagID string) error
}

type TagRepository struct{}

func NewTagRepository() TagRepositoryInterface {
	return &TagRepository{}
}

// CreateTag inserts a new tag record.
// It will set TagID, PK, SK, CreatedAt, UpdatedAt, and default Status if not provided.
func (r *TagRepository) CreateTag(ctx context.Context, db *dynamodb.Client, tag models.Tag) (*models.Tag, error) {
	tagID := uuid.New().String()

	// Partition by CRM so tags can be queried per crm
	tag.PK = fmt.Sprintf("CRM#%s", tag.CrmID)
	tag.SK = fmt.Sprintf("TAG#%s", tagID)

	now := time.Now().UTC()
	tag.CreatedAt = now
	tag.UpdatedAt = now

	if tag.Status == "" {
		tag.Status = "active"
	}

	item, err := attributevalue.MarshalMap(tag)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal tag: %w", err)
	}

	_, err = db.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(TagTableName),
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(PK) AND attribute_not_exists(SK)"),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to put tag: %w", err)
	}

	return &tag, nil
}

// FetchTag retrieves a single tag by crmID + tagID
func (r *TagRepository) FetchTag(ctx context.Context, db *dynamodb.Client, crmID string, tagID string) (*models.Tag, error) {
	pk := fmt.Sprintf("CRM#%s", crmID)
	sk := fmt.Sprintf("TAG#%s", tagID)

	out, err := db.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(TagTableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: sk},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get tag: %w", err)
	}
	if out.Item == nil {
		return nil, fmt.Errorf("tag not found")
	}

	var tag models.Tag
	if err := attributevalue.UnmarshalMap(out.Item, &tag); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tag: %w", err)
	}

	return &tag, nil
}

// ListTagsByCRM queries all tags for a given crmID
func (r *TagRepository) ListTagsByCRM(ctx context.Context, db *dynamodb.Client, crmID string) ([]models.Tag, error) {
	pk := fmt.Sprintf("CRM#%s", crmID)

	out, err := db.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(TagTableName),
		KeyConditionExpression:    aws.String("PK = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{":pk": &types.AttributeValueMemberS{Value: pk}},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query tags: %w", err)
	}

	var tags []models.Tag
	if err := attributevalue.UnmarshalListOfMaps(out.Items, &tags); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tags: %w", err)
	}

	return tags, nil
}

// UpdateTagName updates the tag name and UpdatedAt
func (r *TagRepository) UpdateTagName(ctx context.Context, db *dynamodb.Client, crmID string, tagID string, newName string) error {
	pk := fmt.Sprintf("CRM#%s", crmID)
	sk := fmt.Sprintf("TAG#%s", tagID)
	now := time.Now().UTC()

	updateExpr := "SET #name = :name, #updatedAt = :updatedAt"
	values := map[string]types.AttributeValue{
		":name":      &types.AttributeValueMemberS{Value: newName},
		":updatedAt": &types.AttributeValueMemberS{Value: now.Format(time.RFC3339)},
	}
	names := map[string]string{
		"#name":      "name",
		"#updatedAt": "updatedAt",
	}

	_, err := db.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:                 aws.String(TagTableName),
		Key:                       map[string]types.AttributeValue{"PK": &types.AttributeValueMemberS{Value: pk}, "SK": &types.AttributeValueMemberS{Value: sk}},
		UpdateExpression:          &updateExpr,
		ExpressionAttributeValues: values,
		ExpressionAttributeNames:  names,
	})
	if err != nil {
		return fmt.Errorf("failed to update tag name: %w", err)
	}
	return nil
}

// UpdateTagStatus updates the status and UpdatedAt
func (r *TagRepository) UpdateTagStatus(ctx context.Context, db *dynamodb.Client, crmID string, tagID string, status string) error {
	pk := fmt.Sprintf("CRM#%s", crmID)
	sk := fmt.Sprintf("TAG#%s", tagID)
	now := time.Now().UTC()

	updateExpr := "SET #status = :status, #updatedAt = :updatedAt"
	values := map[string]types.AttributeValue{
		":status":    &types.AttributeValueMemberS{Value: status},
		":updatedAt": &types.AttributeValueMemberS{Value: now.Format(time.RFC3339)},
	}
	names := map[string]string{
		"#status":    "status",
		"#updatedAt": "updatedAt",
	}

	_, err := db.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:                 aws.String(TagTableName),
		Key:                       map[string]types.AttributeValue{"PK": &types.AttributeValueMemberS{Value: pk}, "SK": &types.AttributeValueMemberS{Value: sk}},
		UpdateExpression:          &updateExpr,
		ExpressionAttributeValues: values,
		ExpressionAttributeNames:  names,
	})
	if err != nil {
		return fmt.Errorf("failed to update tag status: %w", err)
	}
	return nil
}

// ArchiveTag marks a tag archived (sets Status and ArchivedAt)
func (r *TagRepository) ArchiveTag(ctx context.Context, db *dynamodb.Client, crmID string, tagID string) error {
	pk := fmt.Sprintf("CRM#%s", crmID)
	sk := fmt.Sprintf("TAG#%s", tagID)
	now := time.Now().UTC()

	updateExpr := "SET #status = :status, #archivedAt = :archivedAt, #updatedAt = :updatedAt"
	values := map[string]types.AttributeValue{
		":status":     &types.AttributeValueMemberS{Value: "archived"},
		":archivedAt": &types.AttributeValueMemberS{Value: now.Format(time.RFC3339)},
		":updatedAt":  &types.AttributeValueMemberS{Value: now.Format(time.RFC3339)},
	}

	names := map[string]string{
		"#status":     "status",
		"#archivedAt": "archivedAt",
		"#updatedAt":  "updatedAt",
	}

	_, err := db.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:                 aws.String(TagTableName),
		Key:                       map[string]types.AttributeValue{"PK": &types.AttributeValueMemberS{Value: pk}, "SK": &types.AttributeValueMemberS{Value: sk}},
		UpdateExpression:          &updateExpr,
		ExpressionAttributeValues: values,
		ExpressionAttributeNames:  names,
	})
	if err != nil {
		return fmt.Errorf("failed to archive tag: %w", err)
	}
	return nil
}
