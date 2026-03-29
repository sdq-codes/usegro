package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type AttributeValue struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Handle string `json:"handle"`
}

type AttributeValues []AttributeValue

func (v AttributeValues) Value() (driver.Value, error) {
	b, err := json.Marshal(v)
	return string(b), err
}

func (v *AttributeValues) Scan(value interface{}) error {
	var b []byte
	switch val := value.(type) {
	case []byte:
		b = val
	case string:
		b = []byte(val)
	default:
		return fmt.Errorf("cannot scan type %T into AttributeValues", value)
	}
	return json.Unmarshal(b, v)
}

type StandardAttribute struct {
	ID         uuid.UUID          `json:"id" gorm:"primaryKey;type:uuid"`
	Name       string             `json:"name" gorm:"type:varchar(255);not null;index"`
	Handle     string             `json:"handle" gorm:"type:varchar(255);not null;uniqueIndex"`
	Values     AttributeValues    `json:"values" gorm:"type:jsonb;not null;default:'[]'"`
	Categories []StandardCategory `json:"-" gorm:"many2many:standard_category_attributes"`
}
