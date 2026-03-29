package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/sdq-codes/usegro-api/internal/apps/form/dtos"
	"github.com/sdq-codes/usegro-api/internal/apps/form/models"
	"github.com/sdq-codes/usegro-api/internal/apps/form/repositories"
	"github.com/sdq-codes/usegro-api/pkg/amplitude"
)

type FormSubmissionService struct {
	formRepo       repositories.FormRepositoryInterface
	submissionRepo repositories.FormSubmissionRepositoryInterface
	activityRepo   repositories.CustomerActivityRepositoryInterface
}

func NewFormSubmissionService(
	formRepo repositories.FormRepositoryInterface,
	submissionRepo repositories.FormSubmissionRepositoryInterface,
	activityRepo repositories.CustomerActivityRepositoryInterface,
) FormSubmissionService {
	return FormSubmissionService{
		formRepo:       formRepo,
		submissionRepo: submissionRepo,
		activityRepo:   activityRepo,
	}
}

func (s *FormSubmissionService) CreateSubmission(
	ctx context.Context,
	formID string,
	versionID string,
	req dtos.CreateSubmissionInput,
	crmID string,
	userID string,
) error {
	form, err := s.formRepo.FetchFormVersion(ctx, formID, versionID)
	if err != nil {
		return err
	}
	if form.Version.FormVersionStatus != "published" {
		return fmt.Errorf("submissions only allowed for published versions")
	}

	email, _ := req.Answers["email"].(string)
	phone, _ := req.Answers["phone_number"].(string)
	dupEmail, dupPhone, err := s.submissionRepo.CheckDuplicateContact(ctx, formID, email, phone)
	if err != nil {
		return fmt.Errorf("failed to check for duplicate contact: %w", err)
	}
	if dupEmail {
		return fmt.Errorf("a customer with this email already exists")
	}
	if dupPhone {
		return fmt.Errorf("a customer with this phone number already exists")
	}

	submissionID := uuid.New().String()
	submission := models.FormSubmission{
		SubmissionID:  submissionID,
		CrmID:         crmID,
		FormID:        formID,
		FormVersionID: versionID,
		Answers:       req.Answers,
		Type:          models.SubmissionType(req.Type),
		VersionSnap:   req.VersionSnap,
	}
	if err := s.submissionRepo.CreateSubmission(ctx, submission); err != nil {
		return err
	}

	amplitude.Track(userID, amplitude.EventContactCreated, map[string]interface{}{
		"crm_id":  crmID,
		"form_id": formID,
		"type":    string(submission.Type),
	})

	if submission.Type == models.SubmissionTypeCustomer && s.activityRepo != nil {
		_ = s.activityRepo.LogActivity(ctx, models.CustomerActivity{
			ActivityType: models.ActivityTypeCustomerCreated,
			Description:  "You created this customer",
			CrmID:        crmID,
			CustomerID:   submissionID,
			FormID:       formID,
			PerformedBy:  userID,
		})
	}
	return nil
}

func (s *FormSubmissionService) ArchiveSubmission(ctx context.Context, formID string, submissionID string) error {
	return s.submissionRepo.UpdateSubmissionStatus(ctx, formID, submissionID, "archived")
}

func (s *FormSubmissionService) UpdateSubmission(ctx context.Context, formID string, submissionID string, answers map[string]interface{}) error {
	return s.submissionRepo.UpdateSubmissionAnswers(ctx, formID, submissionID, answers)
}
