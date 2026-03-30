package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sdq-codes/usegro-api/internal/apps/form/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FormSubmissionRepositoryInterface interface {
	CreateSubmission(ctx context.Context, submission models.FormSubmission) error
	FetchSubmission(ctx context.Context, formID string, submissionID string) (*models.FormSubmission, error)
	UpdateSubmissionStatus(ctx context.Context, formID string, submissionID string, status string) error
	UpdateSubmissionAnswers(ctx context.Context, formID string, submissionID string, answers map[string]interface{}) error
	CheckDuplicateContact(ctx context.Context, formID string, email string, phone string) (duplicateEmail bool, duplicatePhone bool, err error)
}

type FormSubmissionRepository struct {
	db *mongo.Database
}

func NewFormSubmissionRepository(db *mongo.Database) FormSubmissionRepositoryInterface {
	return &FormSubmissionRepository{db: db}
}

func (r *FormSubmissionRepository) col() *mongo.Collection {
	return r.db.Collection("form_submissions")
}

func (r *FormSubmissionRepository) CreateSubmission(ctx context.Context, submission models.FormSubmission) error {
	if submission.SubmissionID == "" {
		submission.SubmissionID = uuid.New().String()
	}
	submission.CreatedAt = time.Now().UTC()
	submission.Status = "active"

	_, err := r.col().InsertOne(ctx, submission)
	if err != nil {
		return fmt.Errorf("failed to create submission: %w", err)
	}
	return nil
}

func (r *FormSubmissionRepository) FetchSubmission(ctx context.Context, formID string, submissionID string) (*models.FormSubmission, error) {
	var sub models.FormSubmission
	err := r.col().FindOne(ctx, bson.M{"_id": submissionID, "formID": formID}).Decode(&sub)
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("submission not found")
	}
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

func (r *FormSubmissionRepository) CheckDuplicateContact(ctx context.Context, formID string, email string, phone string) (duplicateEmail bool, duplicatePhone bool, err error) {
	limit := int64(1)

	if email != "" {
		n, e := r.col().CountDocuments(ctx, bson.M{
			"formID":        formID,
			"answers.email": email,
			"status":        "active",
			"type":          "customer",
		}, &options.CountOptions{Limit: &limit})
		if e != nil {
			err = e
			return
		}
		duplicateEmail = n > 0
	}

	if phone != "" {
		n, e := r.col().CountDocuments(ctx, bson.M{
			"formID":               formID,
			"answers.phone_number": phone,
			"status":               "active",
			"type":                 "customer",
		}, &options.CountOptions{Limit: &limit})
		if e != nil {
			err = e
			return
		}
		duplicatePhone = n > 0
	}
	return
}

func (r *FormSubmissionRepository) UpdateSubmissionStatus(ctx context.Context, formID string, submissionID string, status string) error {
	now := time.Now().UTC()
	_, err := r.col().UpdateOne(ctx,
		bson.M{"_id": submissionID, "formID": formID},
		bson.M{"$set": bson.M{"status": status, "archivedAt": now}},
	)
	return err
}

func (r *FormSubmissionRepository) UpdateSubmissionAnswers(ctx context.Context, formID string, submissionID string, answers map[string]interface{}) error {
	res, err := r.col().UpdateOne(ctx,
		bson.M{"_id": submissionID, "formID": formID},
		bson.M{"$set": bson.M{"answers": answers}},
	)
	if err != nil {
		return fmt.Errorf("failed to update submission answers: %w", err)
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("submission %s not found", submissionID)
	}
	return nil
}

// EnsureSubmissionIndexes creates required MongoDB indexes
func EnsureSubmissionIndexes(ctx context.Context, db *mongo.Database) error {
	_, err := db.Collection("form_submissions").Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "formID", Value: 1}}},
		{Keys: bson.D{{Key: "crmID", Value: 1}, {Key: "type", Value: 1}}},
		{Keys: bson.D{{Key: "formID", Value: 1}, {Key: "answers.email", Value: 1}}},
		{Keys: bson.D{{Key: "formID", Value: 1}, {Key: "answers.phone_number", Value: 1}}},
	})
	return err
}
