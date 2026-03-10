package services

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	formModels "github.com/usegro/services/crm/internal/apps/form/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockCRMCustomerRepository is a mock implementation of the CRM customer repository
type mockCRMCustomerRepository struct {
	mock.Mock
}

func (m *mockCRMCustomerRepository) GetCrmCustomer(ctx context.Context, dynamo *dynamodb.Client, submissionID, formID, crmID string) (*formModels.FormSubmission, error) {
	args := m.Called(ctx, dynamo, submissionID, formID, crmID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*formModels.FormSubmission), args.Error(1)
}

func (m *mockCRMCustomerRepository) FetchCrmCustomers(ctx context.Context, dynamo *dynamodb.Client, crmID string) ([]formModels.FormSubmission, error) {
	args := m.Called(ctx, dynamo, crmID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]formModels.FormSubmission), args.Error(1)
}

func (m *mockCRMCustomerRepository) ArchiveCrmCustomer(ctx context.Context, dynamo *dynamodb.Client, submissionID, formID, crmID string) error {
	args := m.Called(ctx, dynamo, submissionID, formID, crmID)
	return args.Error(0)
}

func (m *mockCRMCustomerRepository) FetchPublishedCreateCustomerForm(ctx context.Context, dynamo *dynamodb.Client, crmID string) (*formModels.CompleteForm, error) {
	args := m.Called(ctx, dynamo, crmID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*formModels.CompleteForm), args.Error(1)
}

func TestGetCrmCustomer_RepositoryError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	submissionID := "submission-123"
	formID := "form-456"
	crmID := "crm-789"
	
	mockRepo := new(mockCRMCustomerRepository)
	expectedError := errors.New("database connection failed")
	
	// Mock repository to return nil customer and an error
	mockRepo.On("GetCrmCustomer", ctx, mock.Anything, submissionID, formID, crmID).
		Return(nil, expectedError)
	
	// Create service with mock repository
	service := &CrmCustomerService{
		crmCustomerRepo: mockRepo,
		dynamo:          nil, // Not needed for this test as we're mocking the repo
	}
	
	// Act
	customer, err := service.GetCrmCustomer(ctx, submissionID, formID, crmID)
	
	// Assert
	assert.Nil(t, customer, "Expected customer to be nil when repository fails")
	assert.Error(t, err, "Expected an error when repository fails")
	assert.Contains(t, err.Error(), "CRM customers could not be deleted", "Expected specific error message")
	
	// Verify mock expectations
	mockRepo.AssertExpectations(t)
}

func TestGetCrmCustomer_Success(t *testing.T) {
	ctx := context.Background()
	submissionID := "submission-123"
	formID := "form-456"
	crmID := "crm-789"
	
	mockRepo := new(mockCRMCustomerRepository)
	expectedCustomer := &formModels.FormSubmission{
		SubmissionID: submissionID,
		FormID:       formID,
		CrmID:        crmID,
	}
	
	mockRepo.On("GetCrmCustomer", ctx, mock.Anything, submissionID, formID, crmID).
		Return(expectedCustomer, nil)
	
	service := &CrmCustomerService{
		crmCustomerRepo: mockRepo,
		dynamo:          nil,
	}
	
	customer, err := service.GetCrmCustomer(ctx, submissionID, formID, crmID)
	
	assert.NoError(t, err)
	assert.NotNil(t, customer)
	assert.Equal(t, expectedCustomer.SubmissionID, customer.SubmissionID)
	assert.Equal(t, expectedCustomer.FormID, customer.FormID)
	assert.Equal(t, expectedCustomer.CrmID, customer.CrmID)
	
	mockRepo.AssertExpectations(t)
}

func TestFetchCrmCustomers_Success(t *testing.T) {
	ctx := context.Background()
	crmID := "crm-789"
	
	mockRepo := new(mockCRMCustomerRepository)
	expectedCustomers := []formModels.FormSubmission{
		{SubmissionID: "sub-1", CrmID: crmID},
		{SubmissionID: "sub-2", CrmID: crmID},
	}
	
	mockRepo.On("FetchCrmCustomers", ctx, mock.Anything, crmID).
		Return(expectedCustomers, nil)
	
	service := &CrmCustomerService{
		crmCustomerRepo: mockRepo,
		dynamo:          nil,
	}
	
	customers, err := service.FetchCrmCustomers(ctx, crmID)
	
	assert.NoError(t, err)
	assert.NotNil(t, customers)
	assert.Len(t, *customers, 2)
	assert.Equal(t, "sub-1", (*customers)[0].SubmissionID)
	
	mockRepo.AssertExpectations(t)
}

func TestFetchCrmCustomers_RepositoryError(t *testing.T) {
	ctx := context.Background()
	crmID := "crm-789"
	
	mockRepo := new(mockCRMCustomerRepository)
	expectedError := errors.New("database error")
	
	mockRepo.On("FetchCrmCustomers", ctx, mock.Anything, crmID).
		Return(nil, expectedError)
	
	service := &CrmCustomerService{
		crmCustomerRepo: mockRepo,
		dynamo:          nil,
	}
	
	customers, err := service.FetchCrmCustomers(ctx, crmID)
	
	assert.Nil(t, customers)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "CRM customers could not be fetched")
	
	mockRepo.AssertExpectations(t)
}

func TestArchiveCrmCustomer_Success(t *testing.T) {
	ctx := context.Background()
	submissionID := "submission-123"
	formID := "form-456"
	crmID := "crm-789"
	
	mockRepo := new(mockCRMCustomerRepository)
	mockRepo.On("ArchiveCrmCustomer", ctx, mock.Anything, submissionID, formID, crmID).
		Return(nil)
	
	service := &CrmCustomerService{
		crmCustomerRepo: mockRepo,
		dynamo:          nil,
	}
	
	err := service.ArchiveCrmCustomer(ctx, submissionID, formID, crmID)
	
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestArchiveCrmCustomer_RepositoryError(t *testing.T) {
	ctx := context.Background()
	submissionID := "submission-123"
	formID := "form-456"
	crmID := "crm-789"
	
	mockRepo := new(mockCRMCustomerRepository)
	expectedError := errors.New("archive failed")
	
	mockRepo.On("ArchiveCrmCustomer", ctx, mock.Anything, submissionID, formID, crmID).
		Return(expectedError)
	
	service := &CrmCustomerService{
		crmCustomerRepo: mockRepo,
		dynamo:          nil,
	}
	
	err := service.ArchiveCrmCustomer(ctx, submissionID, formID, crmID)
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "CRM customers could not be deleted")
	mockRepo.AssertExpectations(t)
}

func TestFetchPublishedCreateCustomerForm_Success(t *testing.T) {
	ctx := context.Background()
	crmID := "crm-789"
	
	mockRepo := new(mockCRMCustomerRepository)
	expectedForm := &formModels.CompleteForm{
		Version: formModels.FormVersion{Title: "Customer Form"},
	}
	
	mockRepo.On("FetchPublishedCreateCustomerForm", ctx, mock.Anything, crmID).
		Return(expectedForm, nil)
	
	service := &CrmCustomerService{
		crmCustomerRepo: mockRepo,
		dynamo:          nil,
	}
	
	form, err := service.FetchPublishedCreateCustomerForm(ctx, crmID)
	
	assert.NoError(t, err)
	assert.NotNil(t, form)
	assert.Equal(t, "Customer Form", form.Version.Title)
	mockRepo.AssertExpectations(t)
}

func TestFetchPublishedCreateCustomerForm_ReturnsFormEvenOnError(t *testing.T) {
	ctx := context.Background()
	crmID := "crm-789"
	
	mockRepo := new(mockCRMCustomerRepository)
	expectedForm := &formModels.CompleteForm{
		Version: formModels.FormVersion{Title: "Form"},
	}
	
	mockRepo.On("FetchPublishedCreateCustomerForm", ctx, mock.Anything, crmID).
		Return(expectedForm, errors.New("some error"))
	
	service := &CrmCustomerService{
		crmCustomerRepo: mockRepo,
		dynamo:          nil,
	}
	
	form, err := service.FetchPublishedCreateCustomerForm(ctx, crmID)
	
	// The service returns the form even if there's an error (based on the implementation)
	assert.NoError(t, err)
	assert.NotNil(t, form)
	mockRepo.AssertExpectations(t)
}
