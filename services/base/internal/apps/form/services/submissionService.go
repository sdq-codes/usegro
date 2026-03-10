package services

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/sdq-codes/usegro-api/internal/apps/form/dtos"
	"github.com/sdq-codes/usegro-api/internal/apps/form/models"
	"github.com/sdq-codes/usegro-api/internal/apps/form/repositories"
)

type FormSubmissionService struct {
	formRepo       repositories.FormRepositoryInterface
	submissionRepo repositories.FormSubmissionRepositoryInterface
	formDynamo     *dynamodb.Client
}

func NewFormSubmissionService(formRepo repositories.FormRepositoryInterface, submissionRepo repositories.FormSubmissionRepositoryInterface, formDynamo *dynamodb.Client) FormSubmissionService {
	return FormSubmissionService{formRepo: formRepo, formDynamo: formDynamo, submissionRepo: submissionRepo}
}

func (s *FormSubmissionService) CreateSubmission(
	ctx context.Context,
	formID string,
	versionID string,
	req dtos.CreateSubmissionInput,
	crmID string,
) error {
	form, err := s.formRepo.FetchFormVersion(ctx, s.formDynamo, formID, versionID)
	if err != nil {
		return err
	}
	if form.Version.FormVersionStatus != "published" {
		return fmt.Errorf("submissions only allowed for published versions")
	}

	/*
		TODO
		@sdq
		Work on form validation for the backend
	*/

	//if err := helpers.ValidateSubmissionAnswers(form.Fields, req.Answers); err != nil {
	//	log.Print(err)
	//	return err
	//}

	submission := models.FormSubmission{
		CrmID:         crmID,
		FormID:        formID,
		FormVersionID: versionID,
		Answers:       req.Answers,
		Type:          models.SubmissionType(req.Type),
		VersionSnap:   req.VersionSnap,
	}
	return s.submissionRepo.CreateSubmission(ctx, s.formDynamo, submission)
}

func (s *FormSubmissionService) ArchiveSubmission(
	ctx context.Context,
	formID string,
	submissionID string,
) error {
	return s.submissionRepo.UpdateSubmissionStatus(ctx, s.formDynamo, formID, submissionID, "archived")
}
