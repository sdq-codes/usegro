package services

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/usegro/services/crm/internal/apps/crm/dto"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/usegro/services/crm/internal/apps/crm/models"
	"github.com/usegro/services/crm/internal/apps/crm/repositories"
	"github.com/usegro/services/crm/pkg/amplitude"
)

type TagService struct {
	tagRepo repositories.TagRepositoryInterface
	dynamo  *dynamodb.Client
}

func NewTagService(dynamo *dynamodb.Client) *TagService {
	return &TagService{
		tagRepo: repositories.NewTagRepository(),
		dynamo:  dynamo,
	}
}

// CreateTag creates a new tag for a CRM
func (s *TagService) CreateTag(ctx context.Context, tag dto.TagCreateDTO, crmId string, userId uuid.UUID) (*models.Tag, error) {
	tagModel := models.Tag{
		CrmID:     crmId,
		CreatedBy: userId.String(),
		Status:    "active",
		Tag:       tag.Tag,
	}
	created, err := s.tagRepo.CreateTag(ctx, s.dynamo, tagModel)
	if err != nil {
		return nil, err
	}
	amplitude.Track(userId.String(), amplitude.EventTagCreated, map[string]interface{}{
		amplitude.PropCrmID:   crmId,
		amplitude.PropTagID:   created.SK,
		amplitude.PropTagName: created.Tag,
	})
	return created, nil
}

// FetchTag retrieves a tag by CRM ID and Tag ID
func (s *TagService) FetchTag(ctx context.Context, crmID string, tagID string) (*models.Tag, error) {
	return s.tagRepo.FetchTag(ctx, s.dynamo, crmID, tagID)
}

// ListTagsByCRM retrieves all tags belonging to a CRM
func (s *TagService) ListTagsByCRM(ctx context.Context, crmID string) ([]models.Tag, error) {
	return s.tagRepo.ListTagsByCRM(ctx, s.dynamo, crmID)
}

// UpdateTagName updates the name of a tag
func (s *TagService) UpdateTagName(ctx context.Context, crmID string, tagID string, newName string) error {
	tag, err := s.tagRepo.FetchTag(ctx, s.dynamo, crmID, tagID)
	if err != nil {
		return err
	}
	if tag == nil {
		return fmt.Errorf("Tag not found")
	}
	return s.tagRepo.UpdateTagName(ctx, s.dynamo, crmID, tagID, newName)
}

// UpdateTagStatus updates the status of a tag
func (s *TagService) UpdateTagStatus(ctx context.Context, crmID string, tagID string) error {
	var status string
	tagName, err := s.tagRepo.FetchTag(ctx, s.dynamo, crmID, tagID)
	if err != nil {
		return err
	}
	if tagName.Status == "active" {
		status = "inactive"
	} else {
		status = "active"
	}
	return s.tagRepo.UpdateTagStatus(ctx, s.dynamo, crmID, tagID, status)
}
