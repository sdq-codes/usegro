package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/usegro/services/crm/internal/apps/crm/dto"
	"github.com/usegro/services/crm/internal/apps/crm/models"
	"github.com/usegro/services/crm/internal/apps/crm/repositories"
	"github.com/usegro/services/crm/pkg/amplitude"
	"go.mongodb.org/mongo-driver/mongo"
)

type TagService struct {
	tagRepo repositories.TagRepositoryInterface
}

func NewTagService(mongoDB *mongo.Database) *TagService {
	return &TagService{
		tagRepo: repositories.NewTagRepository(mongoDB),
	}
}

func (s *TagService) CreateTag(ctx context.Context, tag dto.TagCreateDTO, crmId string, userId uuid.UUID) (*models.Tag, error) {
	tagModel := models.Tag{
		CrmID:     crmId,
		CreatedBy: userId.String(),
		Status:    "active",
		Tag:       tag.Tag,
	}
	created, err := s.tagRepo.CreateTag(ctx, tagModel)
	if err != nil {
		return nil, err
	}
	amplitude.Track(userId.String(), amplitude.EventTagCreated, map[string]interface{}{
		amplitude.PropCrmID:   crmId,
		amplitude.PropTagID:   created.ID,
		amplitude.PropTagName: created.Tag,
	})
	return created, nil
}

func (s *TagService) FetchTag(ctx context.Context, crmID string, tagID string) (*models.Tag, error) {
	return s.tagRepo.FetchTag(ctx, crmID, tagID)
}

func (s *TagService) ListTagsByCRM(ctx context.Context, crmID string) ([]models.Tag, error) {
	return s.tagRepo.ListTagsByCRM(ctx, crmID)
}

func (s *TagService) UpdateTagName(ctx context.Context, crmID string, tagID string, newName string) error {
	tag, err := s.tagRepo.FetchTag(ctx, crmID, tagID)
	if err != nil {
		return err
	}
	if tag == nil {
		return fmt.Errorf("tag not found")
	}
	return s.tagRepo.UpdateTagName(ctx, crmID, tagID, newName)
}

func (s *TagService) UpdateTagStatus(ctx context.Context, crmID string, tagID string) error {
	tag, err := s.tagRepo.FetchTag(ctx, crmID, tagID)
	if err != nil {
		return err
	}
	status := "inactive"
	if tag.Status != "active" {
		status = "active"
	}
	return s.tagRepo.UpdateTagStatus(ctx, crmID, tagID, status)
}
