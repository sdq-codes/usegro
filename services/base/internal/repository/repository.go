package repository

import (
	"github.com/sdq-codes/usegro-api/database"
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
