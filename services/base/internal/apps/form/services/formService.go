package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sdq-codes/usegro-api/internal/apps/form/dtos"
	"github.com/sdq-codes/usegro-api/internal/apps/form/models"
	"github.com/sdq-codes/usegro-api/internal/apps/form/repositories"
	
)

type FormService interface {
	CreateForm(ctx context.Context, req dtos.CreateVersionRequestDTI, crmId string) error
	CreateVersion(ctx context.Context, req dtos.CreateVersionRequestDTI, formId string) error
	CreateFormVersionField(ctx context.Context, formId string, req dtos.CreateVersionFieldInputDTI) error
	FetchForm(ctx context.Context, formID string) (*models.CompleteForm, error)
	FetchDraftForm(ctx context.Context, formID string) (*models.CompleteForm, error)
	FetchFormVersion(ctx context.Context, formID, formVersionID string) (*models.CompleteForm, error)
	PublishFormVersion(ctx context.Context, formId string) error
	DeleteFormVersionField(ctx context.Context, formId string, fieldID string) error
	UpdateFormField(ctx context.Context, formID string, versionID string, fieldID string, updates map[string]interface{}) error
}

type formService struct {
	formRepo repositories.FormRepositoryInterface
}

func NewFormService(formRepo repositories.FormRepositoryInterface) FormService {
	return &formService{formRepo: formRepo}
}

func (s *formService) CreateForm(ctx context.Context, req dtos.CreateVersionRequestDTI, crmId string) error {
	formID := uuid.New().String()
	form := models.Form{
		ID:        formID,
		CrmID:     crmId,
		Type:      req.Type,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	versionID := uuid.New().String()
	formVersion := models.FormVersion{
		ID:                versionID,
		FormID:            formID,
		Title:             req.Title,
		Description:       req.Description,
		FormVersionStatus: "draft",
		CreatedAt:         time.Now().UTC(),
		UpdatedAt:         time.Now().UTC(),
	}

	if err := s.formRepo.CreateForm(ctx, formVersion, form); err != nil {
		return exception2.BadRequestError
	}
	return nil
}

func (s *formService) CreateVersion(ctx context.Context, req dtos.CreateVersionRequestDTI, formId string) error {
	_, err := s.FetchForm(ctx, formId)
	if err != nil {
		return errors.New("form does not exist")
	}
	form, _ := s.FetchDraftForm(ctx, formId)
	if form != nil {
		return errors.New("draft form already exists")
	}

	versionID := uuid.New().String()
	formVersion := models.FormVersion{
		ID:                versionID,
		FormID:            formId,
		Title:             req.Title,
		Description:       req.Description,
		FormVersionStatus: "draft",
		CreatedAt:         time.Now().UTC(),
		UpdatedAt:         time.Now().UTC(),
	}

	if err := s.formRepo.CreateFormVersion(ctx, formVersion); err != nil {
		return exception2.BadRequestError
	}
	return nil
}

func (s *formService) FetchForm(ctx context.Context, formID string) (*models.CompleteForm, error) {
	fullForm, err := s.formRepo.FetchForm(ctx, formID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch form: %w", err)
	}
	return &models.CompleteForm{Version: fullForm.Version, Fields: fullForm.Fields}, nil
}

func (s *formService) FetchDraftForm(ctx context.Context, formID string) (*models.CompleteForm, error) {
	fullForm, err := s.formRepo.FetchDraftForm(ctx, formID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch form: %w", err)
	}
	return &models.CompleteForm{Version: fullForm.Version, Fields: fullForm.Fields}, nil
}

func (s *formService) FetchFormVersion(ctx context.Context, formID, formVersionID string) (*models.CompleteForm, error) {
	return s.formRepo.FetchFormVersion(ctx, formID, formVersionID)
}

func (s *formService) CreateFormVersionField(ctx context.Context, formID string, req dtos.CreateVersionFieldInputDTI) error {
	formVersion, err := s.FetchDraftForm(ctx, formID)
	if err != nil {
		return err
	}
	if formVersion.Version.FormVersionStatus != "draft" {
		return exception2.BadRequestError
	}

	existingForm, err := s.formRepo.FetchForm(ctx, formID)
	if err != nil {
		return fmt.Errorf("failed to fetch fields: %w", err)
	}
	existingFields := existingForm.Fields

	newOrder := req.Order
	if newOrder == 0 {
		newOrder = len(existingFields) + 1
	} else if newOrder > len(existingFields)+1 {
		newOrder = len(existingFields) + 1
	} else {
		for i := range existingFields {
			if existingFields[i].Order >= newOrder {
				existingFields[i].Order++
				if err := s.formRepo.UpdateFormVersionFieldOrder(ctx, existingFields[i]); err != nil {
					return fmt.Errorf("failed to shift field order: %w", err)
				}
			}
		}
	}

	field := models.FormVersionField{
		Label:         req.Label,
		Section:       req.Section,
		FieldTypeID:   req.FieldTypeID,
		FieldTypeName: req.FieldTypeName,
		Configs:       req.Configs,
		Options:       req.Options,
		Validations:   req.Validations,
		Required:      req.Required,
		Placeholder:   req.Placeholder,
		Slug:          req.Slug,
		Description:   req.Description,
		Hint:          req.Hint,
		Order:         newOrder,
	}

	if err := s.formRepo.CreateFormVersionField(ctx, formVersion.Version.ID, field); err != nil {
		return err
	}
	return nil
}

func (s *formService) PublishFormVersion(ctx context.Context, formId string) error {
	formVersion, err := s.FetchDraftForm(ctx, formId)
	if err != nil {
		return err
	}
	return s.formRepo.PublishFormVersion(ctx, formId, formVersion.Version.ID)
}

func (s *formService) DeleteFormVersionField(ctx context.Context, formId string, fieldID string) error {
	formVersion, err := s.FetchDraftForm(ctx, formId)
	if err != nil {
		return err
	}
	if formVersion.Version.FormVersionStatus != "draft" {
		return errors.New("only draft form versions are supported")
	}
	return s.formRepo.DeleteFormVersionField(ctx, formVersion.Version.ID, fieldID)
}

func (s *formService) UpdateFormField(ctx context.Context, formID string, versionID string, fieldID string, updates map[string]interface{}) error {
	form, err := s.formRepo.FetchFormVersion(ctx, formID, versionID)
	if err != nil {
		return err
	}

	var target *models.FormVersionField
	for i := range form.Fields {
		if form.Fields[i].ID == fieldID {
			target = &form.Fields[i]
			break
		}
	}
	if target == nil {
		return fmt.Errorf("field %s not found in form %s version %s", fieldID, formID, versionID)
	}

	if orderVal, ok := updates["order"]; ok {
		switch v := orderVal.(type) {
		case int:
			target.Order = v
		case float64:
			target.Order = int(v)
		default:
			return fmt.Errorf("invalid type for order, got %T", v)
		}
		if err := s.formRepo.UpdateFormVersionFieldOrder(ctx, *target); err != nil {
			return err
		}
		delete(updates, "order")
		delete(updates, "Order")
	}

	if len(updates) > 0 {
		return s.formRepo.UpdateFormVersionField(ctx, fieldID, updates)
	}
	return nil
}
