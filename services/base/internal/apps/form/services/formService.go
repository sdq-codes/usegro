package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamodb2 "github.com/sdq-codes/usegro-api/internal/helper/dynamodb"
	exception2 "github.com/sdq-codes/usegro-api/pkg/exception"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sdq-codes/usegro-api/internal/apps/form/dtos"
	"github.com/sdq-codes/usegro-api/internal/apps/form/models"
	"github.com/sdq-codes/usegro-api/internal/apps/form/repositories"
)

type FormService interface {
	CreateForm(ctx context.Context, req dtos.CreateVersionRequestDTI, crmId string) error
	CreateVersion(ctx context.Context, req dtos.CreateVersionRequestDTI, formId string) error
	CreateFormVersionField(
		ctx context.Context,
		formId string,
		req dtos.CreateVersionFieldInputDTI, // new DTO for a single field
	) error
	FetchForm(ctx context.Context, formID string) (*models.CompleteForm, error)
	FetchDraftForm(ctx context.Context, formID string) (*models.CompleteForm, error)
	FetchFormVersion(ctx context.Context, formID, formVersionID string) (*models.CompleteForm, error)
	PublishFormVersion(ctx context.Context, formId string) error
	DeleteFormVersionField(
		ctx context.Context,
		formId string,
		fieldID string,
	) error
	UpdateFormField(
		ctx context.Context,
		formID string,
		versionID string,
		fieldSK string,
		updates map[string]interface{},
	) error
}

type formService struct {
	formRepo   repositories.FormRepositoryInterface
	formDynamo *dynamodb.Client
}

func NewFormService(formRepo repositories.FormRepositoryInterface, formDynamo *dynamodb.Client) FormService {
	return &formService{formRepo: formRepo, formDynamo: formDynamo}
}

func (s *formService) CreateForm(ctx context.Context, req dtos.CreateVersionRequestDTI, crmId string) error {
	// Step 1: Build the form version model
	form := models.Form{}
	// Generate IDs for Form and initial version
	formID := uuid.New().String()
	form.PK = fmt.Sprintf("FORM#%s", formID)
	form.SK = "METADATA"
	form.CrmID = crmId
	form.Type = req.Type
	form.CreatedAt = time.Now().UTC()
	form.UpdatedAt = form.CreatedAt

	formVersion := models.FormVersion{
		FormID:            formID,
		Title:             req.Title,
		Description:       req.Description,
		FormVersionStatus: "draft", // default new version to draft
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	versionID := uuid.New().String()
	formVersion.PK = fmt.Sprintf("FORM#%s", formVersion.FormID)
	formVersion.SK = fmt.Sprintf("VERSION#%s", versionID)

	// Step 4: Call repository to persist
	if err := s.formRepo.CreateForm(ctx, s.formDynamo, formVersion, form); err != nil {
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
		return errors.New("draft from already exist")
	}
	formVersion := models.FormVersion{
		FormID:            formId,
		Title:             req.Title,
		Description:       req.Description,
		FormVersionStatus: "draft", // default new version to draft
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	versionID := uuid.New().String()
	formVersion.PK = fmt.Sprintf("FORM#%s", formVersion.FormID)
	formVersion.SK = fmt.Sprintf("VERSION#%s", versionID)

	if err := s.formRepo.CreateFormVersion(ctx, s.formDynamo, formVersion); err != nil {
		return exception2.BadRequestError
	}

	return nil
}

func (s *formService) FetchForm(ctx context.Context, formID string) (*models.CompleteForm, error) {
	fullForm, err := s.formRepo.FetchForm(ctx, s.formDynamo, formID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch form: %w", err)
	}

	return &models.CompleteForm{
		Version: fullForm.Version,
		Fields:  fullForm.Fields,
	}, nil
}

func (s *formService) FetchDraftForm(ctx context.Context, formID string) (*models.CompleteForm, error) {
	fullForm, err := s.formRepo.FetchDraftForm(ctx, s.formDynamo, formID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch form: %w", err)
	}

	return &models.CompleteForm{
		Version: fullForm.Version,
		Fields:  fullForm.Fields,
	}, nil
}

func (s *formService) FetchFormVersion(ctx context.Context, formID, formVersionID string) (*models.CompleteForm, error) {
	return s.formRepo.FetchFormVersion(ctx, s.formDynamo, formID, formVersionID)
}

func (s *formService) CreateFormVersionField(
	ctx context.Context,
	formID string,
	req dtos.CreateVersionFieldInputDTI,
) error {
	// Step 1: Fetch draft form version
	formVersion, err := s.FetchDraftForm(ctx, formID)
	if err != nil {
		return err
	}
	if formVersion.Version.FormVersionStatus != "draft" {
		return exception2.BadRequestError
	}

	// Step 2: Fetch existing fields to determine order placement
	existingVersionFields, err := s.formRepo.FetchForm(ctx, s.formDynamo, formID)
	if err != nil {
		return fmt.Errorf("failed to fetch fields: %w", err)
	}
	var existingFields []models.FormVersionField
	existingFields = existingVersionFields.Fields

	newOrder := req.Order
	if newOrder == 0 {
		newOrder = len(existingFields) + 1
	} else if newOrder > len(existingFields)+1 {
		newOrder = len(existingFields) + 1
	} else {
		for i := range existingFields {
			if existingFields[i].Order >= newOrder {
				existingFields[i].Order++
				if err := s.formRepo.UpdateFormVersionFieldOrder(ctx, s.formDynamo, existingFields[i]); err != nil {
					return fmt.Errorf("failed to shift field order: %w", err)
				}
			}
		}
	}
	// Step 4: Create the new field
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
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}
	if err := s.formRepo.CreateFormVersionField(ctx, s.formDynamo, formVersion.Version, field); err != nil {
		return err
	}
	return nil
}

func (s *formService) PublishFormVersion(ctx context.Context, formId string) error {
	formVersion, err := s.FetchDraftForm(ctx, formId)
	if err != nil {
		return err
	}
	return s.formRepo.PublishFormVersion(ctx, s.formDynamo, formId, dynamodb2.ExtractSK(formVersion.Version.SK))
}

func (s *formService) DeleteFormVersionField(
	ctx context.Context,
	formId string,
	fieldID string,
) error {
	formVersion, err := s.FetchDraftForm(ctx, formId)
	if err != nil {
		return err
	}
	if formVersion.Version.FormVersionStatus != "draft" {
		return errors.New("only draft form versions are supported")
	}
	if err := s.formRepo.DeleteFormVersionField(ctx, s.formDynamo, formVersion.Version, fieldID); err != nil {
		return err
	}

	return nil
}

func (s *formService) UpdateFormField(
	ctx context.Context,
	formID string,
	versionID string,
	fieldSK string,
	updates map[string]interface{},
) error {
	// Fetch the form version
	form, err := s.formRepo.FetchFormVersion(ctx, s.formDynamo, formID, versionID)
	if err != nil {
		return err
	}

	// Find the target field
	var target *models.FormVersionField
	for i := range form.Fields {
		if strings.Contains(form.Fields[i].SK, fmt.Sprintf("FIELD#%s", fieldSK)) {
			target = &form.Fields[i]
			break
		}
	}
	if target == nil {
		return fmt.Errorf("field %s not found in form %s version %s", fieldSK, formID, versionID)
	}

	// Special case: order update
	if orderVal, ok := updates["order"]; ok {
		switch v := orderVal.(type) {
		case int:
			target.Order = v
		case float64: // BodyParser into map[string]interface{} often gives numbers as float64
			target.Order = int(v)
		default:
			return fmt.Errorf("invalid type for order, got %T", v)
		}

		// Call specialized repo function
		if err := s.formRepo.UpdateFormVersionFieldOrder(ctx, s.formDynamo, *target); err != nil {
			return err
		}

		// Remove "Order" from updates map so it doesn't get sent to generic updater
		delete(updates, "Order")
	}

	// If there are other properties left to update
	if len(updates) > 0 {
		return s.formRepo.UpdateFormVersionField(ctx, s.formDynamo, *target, updates)
	}

	return nil
}
