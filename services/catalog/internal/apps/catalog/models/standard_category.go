package models

import (
	"github.com/google/uuid"
)

type StandardCategory struct {
	ID         uuid.UUID           `json:"id" gorm:"primaryKey;type:uuid"`
	ParentID   *uuid.UUID          `json:"parent_id,omitempty" gorm:"type:uuid;index"`
	Name       string              `json:"name" gorm:"type:varchar(255);not null;index"`
	FullName   string              `json:"full_name,omitempty" gorm:"type:text"`
	Level      int                 `json:"level" gorm:"default:0"`
	IsLeaf     bool                `json:"is_leaf" gorm:"default:false"`
	Children   []StandardCategory  `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	Attributes []StandardAttribute `json:"attributes,omitempty" gorm:"many2many:standard_category_attributes"`
}
