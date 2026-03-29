package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/usegro/services/crm/internal/apps/crm/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const customerActivityCollection = "customer_activity"

type CustomerActivityRepositoryInterface interface {
	FetchCustomerActivity(ctx context.Context, customerID string) ([]models.CustomerActivity, error)
	LogActivity(ctx context.Context, activity models.CustomerActivity) error
}

type CustomerActivityRepository struct {
	db *mongo.Database
}

func NewCustomerActivityRepository(db *mongo.Database) CustomerActivityRepositoryInterface {
	return &CustomerActivityRepository{db: db}
}

func (r *CustomerActivityRepository) FetchCustomerActivity(ctx context.Context, customerID string) ([]models.CustomerActivity, error) {
	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})
	cur, err := r.db.Collection(customerActivityCollection).Find(ctx, bson.M{"customerID": customerID}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to query customer activity: %w", err)
	}
	defer cur.Close(ctx)

	var activities []models.CustomerActivity
	if err := cur.All(ctx, &activities); err != nil {
		return nil, fmt.Errorf("failed to decode customer activity: %w", err)
	}
	return activities, nil
}

func (r *CustomerActivityRepository) LogActivity(ctx context.Context, activity models.CustomerActivity) error {
	activity.ID = uuid.New().String()
	if activity.CreatedAt.IsZero() {
		activity.CreatedAt = time.Now().UTC()
	}
	_, err := r.db.Collection(customerActivityCollection).InsertOne(ctx, activity)
	if err != nil {
		return fmt.Errorf("failed to log activity: %w", err)
	}
	return nil
}
