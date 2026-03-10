package services

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/sdq-codes/usegro-api/internal/apps/form/dtos"
	"github.com/sdq-codes/usegro-api/internal/apps/form/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockFormRepository is a mock implementation of FormRepositoryInterface
type mockFormRepository struct {
	mock.Mock
}

func (m *mockFormRepository) FetchFormVersion(ctx context.Context, db *dynamodb.Client, formID, versionID string) (*models.CompleteForm, error) {
	args := m.Called(ctx, db, formID, versionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CompleteForm), args.Error(1)
}

func (m *mockFormRepository) CreateForm(ctx context.Context, db *dynamodb.Client, version models.FormVersion, form models.Form) error {
	args := m.Called(ctx, db, version, form)
	return args.Error(0)
}

func (m *mockFormRepository) CreateFormVersion(ctx context.Context, db *dynamodb.Client, version models.FormVersion) error {
	args := m.Called(ctx, db, version)
	return args.Error(0)
}

func (m *mockFormRepository) CreateFormVersionField(ctx context.Context, db *dynamodb.Client, formVersion models.FormVersion, field models.FormVersionField) error {
	args := m.Called(ctx, db, formVersion, field)
	return args.Error(0)
}

func (m *mockFormRepository) UpdateFormVersionFieldOrder(ctx context.Context, db *dynamodb.Client, field models.FormVersionField) error {
	args := m.Called(ctx, db, field)
	return args.Error(0)
}

func (m *mockFormRepository) FetchForm(ctx context.Context, db *dynamodb.Client, formID string) (*models.CompleteForm, error) {
	args := m.Called(ctx, db, formID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CompleteForm), args.Error(1)
}

func (m *mockFormRepository) FetchDraftForm(ctx context.Context, db *dynamodb.Client, formID string) (*models.CompleteForm, error) {
	args := m.Called(ctx, db, formID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CompleteForm), args.Error(1)
}

func (m *mockFormRepository) PublishFormVersion(ctx context.Context, db *dynamodb.Client, formID, versionID string) error {
	args := m.Called(ctx, db, formID, versionID)
	return args.Error(0)
}

func (m *mockFormRepository) DeleteFormVersionField(ctx context.Context, db *dynamodb.Client, formVersion models.FormVersion, fieldID string) error {
	args := m.Called(ctx, db, formVersion, fieldID)
	return args.Error(0)
}

func (m *mockFormRepository) UpdateFormVersionField(ctx context.Context, db *dynamodb.Client, field models.FormVersionField, updates map[string]interface{}) error {
	args := m.Called(ctx, db, field, updates)
	return args.Error(0)
}

// mockFormSubmissionRepository is a mock implementation of FormSubmissionRepositoryInterface
type mockFormSubmissionRepository struct {
	mock.Mock
}

func (m *mockFormSubmissionRepository) CreateSubmission(ctx context.Context, db *dynamodb.Client, submission models.FormSubmission) error {
	args := m.Called(ctx, db, submission)
	return args.Error(0)
}

func (m *mockFormSubmissionRepository) FetchSubmission(ctx context.Context, db *dynamodb.Client, formID, submissionID string) (*models.FormSubmission, error) {
	args := m.Called(ctx, db, formID, submissionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.FormSubmission), args.Error(1)
}

func (m *mockFormSubmissionRepository) UpdateSubmissionStatus(ctx context.Context, db *dynamodb.Client, formID, submissionID, status string) error {
	args := m.Called(ctx, db, formID, submissionID, status)
	return args.Error(0)
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

	mockFormRepo.On("FetchFormVersion", ctx, mock.Anything, formID, versionID).
		Return(completeForm, nil)

	mockSubRepo.On("CreateSubmission", ctx, mock.Anything, mock.MatchedBy(func(sub models.FormSubmission) bool {
		return sub.CrmID == crmID && sub.FormID == formID && sub.FormVersionID == versionID
	})).Return(nil)

	service := NewFormSubmissionService(mockFormRepo, mockSubRepo, nil)
	err := service.CreateSubmission(ctx, formID, versionID, submissionInput, crmID)

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
	mockFormRepo.On("FetchFormVersion", ctx, mock.Anything, formID, versionID).
		Return(nil, expectedError)

	submissionInput := dtos.CreateSubmissionInput{
		Answers: map[string]interface{}{},
		Type:    "customer",
	}

	service := NewFormSubmissionService(mockFormRepo, mockSubRepo, nil)
	err := service.CreateSubmission(ctx, formID, versionID, submissionInput, crmID)

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

	mockFormRepo.On("FetchFormVersion", ctx, mock.Anything, formID, versionID).
		Return(completeForm, nil)

	submissionInput := dtos.CreateSubmissionInput{
		Answers: map[string]interface{}{},
		Type:    "customer",
	}

	service := NewFormSubmissionService(mockFormRepo, mockSubRepo, nil)
	err := service.CreateSubmission(ctx, formID, versionID, submissionInput, crmID)

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

	mockFormRepo.On("FetchFormVersion", ctx, mock.Anything, formID, versionID).
		Return(completeForm, nil)

	mockSubRepo.On("CreateSubmission", ctx, mock.Anything, mock.Anything).
		Return(expectedError)

	service := NewFormSubmissionService(mockFormRepo, mockSubRepo, nil)
	err := service.CreateSubmission(ctx, formID, versionID, submissionInput, crmID)

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

	mockSubRepo.On("UpdateSubmissionStatus", ctx, mock.Anything, formID, submissionID, "archived").
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

	mockSubRepo.On("UpdateSubmissionStatus", ctx, mock.Anything, formID, submissionID, "archived").
		Return(expectedError)

	service := NewFormSubmissionService(mockFormRepo, mockSubRepo, nil)
	err := service.ArchiveSubmission(ctx, formID, submissionID)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockSubRepo.AssertExpectations(t)
}
