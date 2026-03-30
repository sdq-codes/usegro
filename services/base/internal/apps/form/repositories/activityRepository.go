package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sdq-codes/usegro-api/internal/apps/form/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const customerActivityCollection = "customer_activity"

type CustomerActivityRepositoryInterface interface {
	LogActivity(ctx context.Context, activity models.CustomerActivity) error
}

type CustomerActivityRepository struct {
	db *mongo.Database
}

func NewCustomerActivityRepository(db *mongo.Database) CustomerActivityRepositoryInterface {
	return &CustomerActivityRepository{db: db}
}

func (r *CustomerActivityRepository) LogActivity(ctx context.Context, activity models.CustomerActivity) error {
	activity.ID = uuid.New().String()
	activity.CreatedAt = time.Now().UTC()

	_, err := r.db.Collection(customerActivityCollection).InsertOne(ctx, activity)
	return err
}

// EnsureActivityIndexes creates required indexes
func EnsureActivityIndexes(ctx context.Context, db *mongo.Database) error {
	_, err := db.Collection(customerActivityCollection).Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "customerID", Value: 1}, {Key: "createdAt", Value: -1}}},
	})
	return err
}
