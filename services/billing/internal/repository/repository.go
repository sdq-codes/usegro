package repository

import (
	"github.com/usegro/services/billing/database"
	"gorm.io/gorm"
)

type Repository struct {
	Db *gorm.DB
}

func NewRepository() *Repository {
	return &Repository{
		Db: database.PostgressInstance1,
	}
}
