package repositories

import (
	"context"
	"github.com/google/uuid"
	"github.com/sdq-codes/usegro-api/internal/apps/base/models"
	"gorm.io/gorm"
)

type JobRepositoryInterface interface {
	AddJob(ctx context.Context, tx *gorm.DB, job *models.Job) (*models.Job, error)
	UpdateJobStatus(ctx context.Context, id uuid.UUID, processing string) error
	AddFailedJob(ctx context.Context, job models.FailedJob) (failedJobID int, err error)
	ResetProcessingJobsToPending(ctx context.Context) error
	GetJobs(ctx context.Context) ([]models.Job, error)
	GetUnfinishedJobs(ctx context.Context) ([]models.Job, error)
	GetFailedJobs(ctx context.Context) ([]models.FailedJob, error)
	RemoveFailedJob(ctx context.Context, jobID uuid.UUID) error
	FetchJob(ctx context.Context, tx *gorm.DB, jobId string) (*models.Job, error)
}

type JobRepository struct {
	db *gorm.DB
}

func NewJobRepository(db *gorm.DB) JobRepositoryInterface {
	return &JobRepository{db: db}
}

func (j *JobRepository) FetchJob(ctx context.Context, tx *gorm.DB, jobId string) (*models.Job, error) {
	job := &models.Job{}
	// WHERE must come before First — First() executes the query immediately in GORM v2.
	if err := tx.WithContext(ctx).Where("id = ?", jobId).First(job).Error; err != nil {
		return nil, err
	}
	return job, nil
}

func (j *JobRepository) AddJob(ctx context.Context, tx *gorm.DB, job *models.Job) (*models.Job, error) {
	job.ID = uuid.New()
	if err := tx.WithContext(ctx).Create(job).Error; err != nil {
		return nil, err
	}
	return job, nil
}

func (j *JobRepository) UpdateJobStatus(ctx context.Context, id uuid.UUID, processing string) error {
	panic("implement me")
}

func (j *JobRepository) AddFailedJob(ctx context.Context, job models.FailedJob) (failedJobID int, err error) {
	panic("implement me")
}

func (j *JobRepository) ResetProcessingJobsToPending(ctx context.Context) error {
	panic("implement me")
}

func (j *JobRepository) GetJobs(ctx context.Context) ([]models.Job, error) {
	panic("implement me")
}

func (j *JobRepository) GetUnfinishedJobs(ctx context.Context) ([]models.Job, error) {
	panic("implement me")
}

func (j *JobRepository) GetFailedJobs(ctx context.Context) ([]models.FailedJob, error) {
	panic("implement me")
}

func (j *JobRepository) RemoveFailedJob(ctx context.Context, jobID uuid.UUID) error {
	panic("implement me")
}
