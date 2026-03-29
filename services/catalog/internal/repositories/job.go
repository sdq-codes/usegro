package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/usegro/services/catalog/internal/models"
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
	return j.db.WithContext(ctx).Model(&models.Job{}).Where("id = ?", id).Update("status", processing).Error
}

func (j *JobRepository) AddFailedJob(ctx context.Context, job models.FailedJob) (failedJobID int, err error) {
	if err := j.db.WithContext(ctx).Create(&job).Error; err != nil {
		return 0, err
	}
	return job.ID, nil
}

func (j *JobRepository) ResetProcessingJobsToPending(ctx context.Context) error {
	return j.db.WithContext(ctx).Model(&models.Job{}).Where("status = ?", "processing").Update("status", "pending").Error
}

func (j *JobRepository) GetJobs(ctx context.Context) ([]models.Job, error) {
	var jobs []models.Job
	if err := j.db.WithContext(ctx).Find(&jobs).Error; err != nil {
		return nil, err
	}
	return jobs, nil
}

func (j *JobRepository) GetUnfinishedJobs(ctx context.Context) ([]models.Job, error) {
	var jobs []models.Job
	if err := j.db.WithContext(ctx).Where("status IN ?", []string{"pending", "processing"}).Find(&jobs).Error; err != nil {
		return nil, err
	}
	return jobs, nil
}

func (j *JobRepository) GetFailedJobs(ctx context.Context) ([]models.FailedJob, error) {
	var jobs []models.FailedJob
	if err := j.db.WithContext(ctx).Find(&jobs).Error; err != nil {
		return nil, err
	}
	return jobs, nil
}

func (j *JobRepository) RemoveFailedJob(ctx context.Context, jobID uuid.UUID) error {
	return j.db.WithContext(ctx).Where("job_id = ?", jobID).Delete(&models.FailedJob{}).Error
}
