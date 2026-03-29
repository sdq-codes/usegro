package services

import (
	"context"
	"errors"
	"testing"

	"github.com/sdq-codes/usegro-api/internal/apps/form/dtos"
	"github.com/sdq-codes/usegro-api/internal/apps/form/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockFormRepository is a mock implementation of FormRepositoryInterface
type mockFormRepository struct {
	mock.Mock
}

func (m *mockFormRepository) FetchFormVersion(ctx context.Context, formID, versionID string) (*models.CompleteForm, error) {
	args := m.Called(ctx, formID, versionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CompleteForm), args.Error(1)
}

func (m *mockFormRepository) CreateForm(ctx context.Context, version models.FormVersion, form models.Form) error {
	args := m.Called(ctx, version, form)
	return args.Error(0)
}

func (m *mockFormRepository) CreateFormVersion(ctx context.Context, version models.FormVersion) error {
	args := m.Called(ctx, version)
	return args.Error(0)
}

func (m *mockFormRepository) CreateFormVersionField(ctx context.Context, formVersionID string, field models.FormVersionField) error {
	args := m.Called(ctx, formVersionID, field)
	return args.Error(0)
}

func (m *mockFormRepository) UpdateFormVersionFieldOrder(ctx context.Context, field models.FormVersionField) error {
	args := m.Called(ctx, field)
	return args.Error(0)
}

func (m *mockFormRepository) FetchForm(ctx context.Context, formID string) (*models.CompleteForm, error) {
	args := m.Called(ctx, formID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CompleteForm), args.Error(1)
}

func (m *mockFormRepository) FetchDraftForm(ctx context.Context, formID string) (*models.CompleteForm, error) {
	args := m.Called(ctx, formID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CompleteForm), args.Error(1)
}

func (m *mockFormRepository) PublishFormVersion(ctx context.Context, formID, versionID string) error {
	args := m.Called(ctx, formID, versionID)
	return args.Error(0)
}

func (m *mockFormRepository) DeleteFormVersionField(ctx context.Context, formVersionID string, fieldID string) error {
	args := m.Called(ctx, formVersionID, fieldID)
	return args.Error(0)
}

func (m *mockFormRepository) UpdateFormVersionField(ctx context.Context, fieldID string, updates map[string]interface{}) error {
	args := m.Called(ctx, fieldID, updates)
	return args.Error(0)
}

// mockFormSubmissionRepository is a mock implementation of FormSubmissionRepositoryInterface
type mockFormSubmissionRepository struct {
	mock.Mock
}

func (m *mockFormSubmissionRepository) CreateSubmission(ctx context.Context, submission models.FormSubmission) error {
	args := m.Called(ctx, submission)
	return args.Error(0)
}

func (m *mockFormSubmissionRepository) FetchSubmission(ctx context.Context, formID, submissionID string) (*models.FormSubmission, error) {
	args := m.Called(ctx, formID, submissionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.FormSubmission), args.Error(1)
}

func (m *mockFormSubmissionRepository) UpdateSubmissionStatus(ctx context.Context, formID, submissionID, status string) error {
	args := m.Called(ctx, formID, submissionID, status)
	return args.Error(0)
}

func (m *mockFormSubmissionRepository) CheckDuplicateContact(ctx context.Context, formID, email, phone string) (bool, bool, error) {
	args := m.Called(ctx, formID, email, phone)
	return args.Bool(0), args.Bool(1), args.Error(2)
}

func TestCreateSubmission_Success(t *testing.T) {
	ctx := context.Background()
	formID := "form-123"
	versionID := "version-456"
	crmID := "crm-789"

	mockFormRepo := new(mockFormRepository)
	mockSubRepo := new(mockFormSubmissionRepository)

	completeForm := &models.CompleteForm{
		Version: models.FormVersion{
			FormVersionStatus: "published",
		},
	}

	submissionInput := dtos.CreateSubmissionInput{
		Answers: map[string]interface{}{
			"field1": "value1",
		},
		Type:        "customer",
		VersionSnap: map[string]interface{}{},
	}

	mockFormRepo.On("FetchFormVersion", ctx, formID, versionID).
		Return(completeForm, nil)

	mockSubRepo.On("CheckDuplicateContact", ctx, formID, "", "").
		Return(false, false, nil)

	mockSubRepo.On("CreateSubmission", ctx, mock.MatchedBy(func(sub models.FormSubmission) bool {
		return sub.CrmID == crmID && sub.FormID == formID && sub.FormVersionID == versionID
	})).Return(nil)

	service := NewFormSubmissionService(mockFormRepo, mockSubRepo, nil)
	err := service.CreateSubmission(ctx, formID, versionID, submissionInput, crmID, "")

	assert.NoError(t, err)
	mockFormRepo.AssertExpectations(t)
	mockSubRepo.AssertExpectations(t)
}

func TestCreateSubmission_FormVersionNotFound(t *testing.T) {
	ctx := context.Background()
	formID := "form-123"
	versionID := "version-456"
	crmID := "crm-789"

	mockFormRepo := new(mockFormRepository)
	mockSubRepo := new(mockFormSubmissionRepository)

	expectedError := errors.New("form version not found")
	mockFormRepo.On("FetchFormVersion", ctx, formID, versionID).
		Return(nil, expectedError)

	submissionInput := dtos.CreateSubmissionInput{
		Answers: map[string]interface{}{},
		Type:    "customer",
	}

	service := NewFormSubmissionService(mockFormRepo, mockSubRepo, nil)
	err := service.CreateSubmission(ctx, formID, versionID, submissionInput, crmID, "")

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockFormRepo.AssertExpectations(t)
	mockSubRepo.AssertNotCalled(t, "CreateSubmission")
}

func TestCreateSubmission_VersionNotPublished(t *testing.T) {
	ctx := context.Background()
	formID := "form-123"
	versionID := "version-456"
	crmID := "crm-789"

	mockFormRepo := new(mockFormRepository)
	mockSubRepo := new(mockFormSubmissionRepository)

	completeForm := &models.CompleteForm{
		Version: models.FormVersion{
			FormVersionStatus: "draft",
		},
	}

	mockFormRepo.On("FetchFormVersion", ctx, formID, versionID).
		Return(completeForm, nil)

	submissionInput := dtos.CreateSubmissionInput{
		Answers: map[string]interface{}{},
		Type:    "customer",
	}

	service := NewFormSubmissionService(mockFormRepo, mockSubRepo, nil)
	err := service.CreateSubmission(ctx, formID, versionID, submissionInput, crmID, "")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "submissions only allowed for published versions")
	mockFormRepo.AssertExpectations(t)
	mockSubRepo.AssertNotCalled(t, "CreateSubmission")
}

func TestCreateSubmission_RepositoryError(t *testing.T) {
	ctx := context.Background()
	formID := "form-123"
	versionID := "version-456"
	crmID := "crm-789"

	mockFormRepo := new(mockFormRepository)
	mockSubRepo := new(mockFormSubmissionRepository)

	completeForm := &models.CompleteForm{
		Version: models.FormVersion{
			FormVersionStatus: "published",
		},
	}

	submissionInput := dtos.CreateSubmissionInput{
		Answers: map[string]interface{}{
			"field1": "value1",
		},
		Type: "customer",
	}

	expectedError := errors.New("database error")

	mockFormRepo.On("FetchFormVersion", ctx, formID, versionID).
		Return(completeForm, nil)

	mockSubRepo.On("CheckDuplicateContact", ctx, formID, "", "").
		Return(false, false, nil)

	mockSubRepo.On("CreateSubmission", ctx, mock.Anything).
		Return(expectedError)

	service := NewFormSubmissionService(mockFormRepo, mockSubRepo, nil)
	err := service.CreateSubmission(ctx, formID, versionID, submissionInput, crmID, "")

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockFormRepo.AssertExpectations(t)
	mockSubRepo.AssertExpectations(t)
}

func TestArchiveSubmission_Success(t *testing.T) {
	ctx := context.Background()
	formID := "form-123"
	submissionID := "submission-456"

	mockFormRepo := new(mockFormRepository)
	mockSubRepo := new(mockFormSubmissionRepository)

	mockSubRepo.On("UpdateSubmissionStatus", ctx, formID, submissionID, "archived").
		Return(nil)

	service := NewFormSubmissionService(mockFormRepo, mockSubRepo, nil)
	err := service.ArchiveSubmission(ctx, formID, submissionID)

	assert.NoError(t, err)
	mockSubRepo.AssertExpectations(t)
}

func TestArchiveSubmission_RepositoryError(t *testing.T) {
	ctx := context.Background()
	formID := "form-123"
	submissionID := "submission-456"

	mockFormRepo := new(mockFormRepository)
	mockSubRepo := new(mockFormSubmissionRepository)

	expectedError := errors.New("update failed")

	mockSubRepo.On("UpdateSubmissionStatus", ctx, formID, submissionID, "archived").
		Return(expectedError)

	service := NewFormSubmissionService(mockFormRepo, mockSubRepo, nil)
	err := service.ArchiveSubmission(ctx, formID, submissionID)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockSubRepo.AssertExpectations(t)
}
