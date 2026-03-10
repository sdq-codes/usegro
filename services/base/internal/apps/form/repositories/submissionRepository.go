package repositories

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/sdq-codes/usegro-api/internal/apps/form/models"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

type FormSubmissionRepositoryInterface interface {
	CreateSubmission(
		ctx context.Context,
		db *dynamodb.Client,
		submission models.FormSubmission,
	) error
	FetchSubmission(
		ctx context.Context,
		db *dynamodb.Client,
		formID string,
		submissionID string,
	) (*models.FormSubmission, error)
	UpdateSubmissionStatus(
		ctx context.Context,
		db *dynamodb.Client,
		formID string,
		submissionID string,
		status string,
	) error
}

type FormSubmissionRepository struct{}

func NewFormSubmissionRepository() FormSubmissionRepositoryInterface {
	return &FormSubmissionRepository{}
}

func (r *FormSubmissionRepository) CreateSubmission(
	ctx context.Context,
	db *dynamodb.Client,
	submission models.FormSubmission,
) error {
	submission.SubmissionID = uuid.New().String()
	submission.PK = fmt.Sprintf("FORM#%s", submission.FormID)
	submission.SK = fmt.Sprintf("SUBMISSION#%s", submission.SubmissionID)
	submission.CreatedAt = time.Now().UTC()
	submission.Status = "active"
	item, err := attributevalue.MarshalMap(submission)
	if err != nil {
		return fmt.Errorf("failed to marshal submission: %w", err)
	}
	_, err = db.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String("form_submissions"),
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(PK) AND attribute_not_exists(SK)"),
	})
	return err
}

func (r *FormSubmissionRepository) FetchSubmission(
	ctx context.Context,
	db *dynamodb.Client,
	formID string,
	submissionID string,
) (*models.FormSubmission, error) {
	pk := fmt.Sprintf("FORM#%s", formID)
	sk := fmt.Sprintf("SUBMISSION#%s", submissionID)

	out, err := db.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String("form_submissions"),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: sk},
		},
	})
	if err != nil {
		return nil, err
	}
	if out.Item == nil {
		return nil, fmt.Errorf("submission not found")
	}

	var sub models.FormSubmission
	if err := attributevalue.UnmarshalMap(out.Item, &sub); err != nil {
		return nil, err
	}
	return &sub, nil
}

func (r *FormSubmissionRepository) UpdateSubmissionStatus(
	ctx context.Context,
	db *dynamodb.Client,
	formID string,
	submissionID string,
	status string,
) error {
	pk := fmt.Sprintf("FORM#%s", formID)
	sk := fmt.Sprintf("SUBMISSION#%s", submissionID)
	now := time.Now().UTC()

	updateExpr := "SET #status = :status, #archivedAt = :archivedAt"
	values := map[string]types.AttributeValue{
		":status":     &types.AttributeValueMemberS{Value: status},
		":archivedAt": &types.AttributeValueMemberS{Value: now.Format(time.RFC3339)},
	}
	names := map[string]string{
		"#status":     "Status",
		"#archivedAt": "ArchivedAt",
	}

	_, err := db.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String("form_submissions"),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: sk},
		},
		UpdateExpression:          &updateExpr,
		ExpressionAttributeValues: values,
		ExpressionAttributeNames:  names,
	})
	return err
}
