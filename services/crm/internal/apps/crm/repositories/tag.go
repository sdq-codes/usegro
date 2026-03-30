package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/usegro/services/crm/internal/apps/crm/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const tagCollection = "tags"

type TagRepositoryInterface interface {
	CreateTag(ctx context.Context, tag models.Tag) (*models.Tag, error)
	FetchTag(ctx context.Context, crmID string, tagID string) (*models.Tag, error)
	ListTagsByCRM(ctx context.Context, crmID string) ([]models.Tag, error)
	UpdateTagName(ctx context.Context, crmID string, tagID string, newName string) error
	UpdateTagStatus(ctx context.Context, crmID string, tagID string, status string) error
	ArchiveTag(ctx context.Context, crmID string, tagID string) error
}

type TagRepository struct {
	db *mongo.Database
}

func NewTagRepository(db *mongo.Database) TagRepositoryInterface {
	return &TagRepository{db: db}
}

func (r *TagRepository) col() *mongo.Collection { return r.db.Collection(tagCollection) }

func (r *TagRepository) CreateTag(ctx context.Context, tag models.Tag) (*models.Tag, error) {
	tag.ID = uuid.New().String()
	now := time.Now().UTC()
	tag.CreatedAt = now
	tag.UpdatedAt = now
	if tag.Status == "" {
		tag.Status = "active"
	}

	_, err := r.col().InsertOne(ctx, tag)
	if err != nil {
		return nil, fmt.Errorf("failed to create tag: %w", err)
	}
	return &tag, nil
}

func (r *TagRepository) FetchTag(ctx context.Context, crmID string, tagID string) (*models.Tag, error) {
	var tag models.Tag
	err := r.col().FindOne(ctx, bson.M{"_id": tagID, "crmID": crmID}).Decode(&tag)
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("tag not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get tag: %w", err)
	}
	return &tag, nil
}

func (r *TagRepository) ListTagsByCRM(ctx context.Context, crmID string) ([]models.Tag, error) {
	cur, err := r.col().Find(ctx, bson.M{"crmID": crmID})
	if err != nil {
		return nil, fmt.Errorf("failed to query tags: %w", err)
	}
	defer cur.Close(ctx)

	var tags []models.Tag
	if err := cur.All(ctx, &tags); err != nil {
		return nil, fmt.Errorf("failed to decode tags: %w", err)
	}
	return tags, nil
}

func (r *TagRepository) UpdateTagName(ctx context.Context, crmID string, tagID string, newName string) error {
	_, err := r.col().UpdateOne(ctx,
		bson.M{"_id": tagID, "crmID": crmID},
		bson.M{"$set": bson.M{"tag": newName, "updatedAt": time.Now().UTC()}},
	)
	return err
}

func (r *TagRepository) UpdateTagStatus(ctx context.Context, crmID string, tagID string, status string) error {
	_, err := r.col().UpdateOne(ctx,
		bson.M{"_id": tagID, "crmID": crmID},
		bson.M{"$set": bson.M{"status": status, "updatedAt": time.Now().UTC()}},
	)
	return err
}

func (r *TagRepository) ArchiveTag(ctx context.Context, crmID string, tagID string) error {
	now := time.Now().UTC()
	_, err := r.col().UpdateOne(ctx,
		bson.M{"_id": tagID, "crmID": crmID},
		bson.M{"$set": bson.M{"status": "archived", "archivedAt": now, "updatedAt": now}},
	)
	return err
}
