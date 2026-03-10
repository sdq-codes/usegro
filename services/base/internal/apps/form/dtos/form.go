package dtos

import "github.com/sdq-codes/usegro-api/internal/apps/form/models"

type CreateVersionRequestDTI struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Type        string `json:"type"`
}
type CreateVersionFieldInputDTI struct {
	Label         string              `json:"label"`
	Description   string              `json:"description"`
	Section       string              `json:"section"`
	FieldTypeID   uint                `json:"fieldTypeID"`
	FieldTypeName string              `json:"fieldTypeName"`
	Options       []models.Option     `json:"options"`
	Configs       []map[string]string `json:"configs"`
	Validations   []map[string]string `json:"validations"`
	Required      bool                `json:"required"`
	Slug          string              `json:"slug"`
	Hint          string              `json:"hint"`
	Logic         []FieldLogicDTO     `json:"logic"`
	Order         int                 `json:"order"`
	Placeholder   string              `json:"placeholder"`
}

type FieldLogicDTO struct {
	FormVersionFieldID string      `json:"formVersionFieldID" validate:"required"` // the field it depends on
	Operator           string      `json:"operator" validate:"required"`           // e.g. equals, not_equals, greater_than
	Value              interface{} `json:"value" validate:"required"`              // the value to compare
	Action             string      `json:"action" validate:"required"`             // e.g. show, hide, enable, disable
}
