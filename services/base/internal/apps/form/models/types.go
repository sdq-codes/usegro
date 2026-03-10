package models

import (
	"github.com/sdq-codes/usegro-api/internal/apps/form/seeds"
	"gorm.io/gorm"
	"time"
)

type FieldType struct {
	ID          uint                  `gorm:"primaryKey" json:"id"`
	Name        string                `gorm:"uniqueIndex;not null" json:"name"` // e.g. "Short Text", "Dropdown"
	Description string                `gorm:"type:text" json:"description"`
	Configs     []FieldTypeConfig     `gorm:"foreignKey:FieldTypeID;constraint:OnDelete:CASCADE" json:"configs"`
	Validations []FieldTypeValidation `gorm:"foreignKey:FieldTypeID;constraint:OnDelete:CASCADE" json:"validations"`
	CreatedAt   time.Time             `gorm:"default:now()" json:"created_at"`
	UpdatedAt   time.Time             `gorm:"default:now()" json:"updated_at"`
}

type FieldTypeConfig struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	FieldTypeID uint      `gorm:"index;not null" json:"field_type_id"`
	Key         string    `gorm:"not null" json:"key"`        // e.g. "placeholder", "minLength", "maxLength"
	ValueType   string    `gorm:"not null" json:"value_type"` // e.g. "string", "int", "boolean"
	Description string    `gorm:"type:text;default:null" json:"description"`
	CreatedAt   time.Time `gorm:"default:now()" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:now()" json:"updated_at"`
}

type FieldTypeValidation struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	FieldTypeID uint      `gorm:"index;not null" json:"field_type_id"`
	Key         string    `gorm:"not null" json:"key"`     // e.g. "required", "regex", "min", "max"
	Rule        string    `gorm:"not null" json:"rule"`    // e.g. regex pattern, int, boolean
	Message     string    `gorm:"not null" json:"message"` // e.g. "This field is required"
	CreatedAt   time.Time `gorm:"default:now()" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:now()" json:"updated_at"`
}

func SeedFieldTypes(db *gorm.DB) error {
	for _, f := range seeds.FieldTypesSeed {
		var fieldType FieldType
		err := db.Where("name = ?", f["name"]).First(&fieldType).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}

		// If not found, create
		if fieldType.ID == 0 {
			fieldType = FieldType{
				Name:        f["name"].(string),
				Description: f["description"].(string),
			}
			if err := db.Create(&fieldType).Error; err != nil {
				return err
			}
		} else {
			// Update description if changed
			fieldType.Description = f["description"].(string)
			if err := db.Save(&fieldType).Error; err != nil {
				return err
			}
		}

		// Sync configs
		for _, c := range f["configs"].([]map[string]string) {
			var config FieldTypeConfig
			if err := db.Where("field_type_id = ? AND key = ?", fieldType.ID, c["key"]).First(&config).Error; err == gorm.ErrRecordNotFound {
				config = FieldTypeConfig{
					FieldTypeID: fieldType.ID,
					Key:         c["key"],
					ValueType:   c["valueType"],
					Description: c["description"],
				}
				db.Create(&config)
			} else {
				config.ValueType = c["valueType"]
				config.Description = c["description"]
				db.Save(&config)
			}
		}

		// Sync validations
		for _, v := range f["validations"].([]map[string]string) {
			var validation FieldTypeValidation
			if err := db.Where("field_type_id = ? AND key = ?", fieldType.ID, v["key"]).First(&validation).Error; err == gorm.ErrRecordNotFound {
				validation = FieldTypeValidation{
					FieldTypeID: fieldType.ID,
					Key:         v["key"],
					Rule:        v["rule"],
					Message:     v["message"],
				}
				db.Create(&validation)
			} else {
				validation.Rule = v["rule"]
				validation.Message = v["message"]
				db.Save(&validation)
			}
		}
	}
	return nil
}
