package models

import "time"

type Option struct {
	Label string `dynamodbav:"label" json:"label"`
	Value string `dynamodbav:"value" json:"value"`
}

type Alert struct {
	Icon    string `dynamodbav:"icon" json:"icon"`
	Type    string `dynamodbav:"type" json:"type"`
	Message string `dynamodbav:"message" json:"message"`
}

type FormVersionField struct {
	PK            string              `dynamodbav:"PK"`
	SK            string              `dynamodbav:"SK"`
	FormVersionID string              `dynamodbav:"formVersionID" json:"formVersionID"`
	FieldTypeID   uint                `dynamodbav:"fieldTypeId" json:"fieldTypeID"`
	FieldTypeName string              `dynamodbav:"fieldTypeName" json:"fieldTypeName"`
	Label         string              `dynamodbav:"label" json:"label"`
	Description   string              `dynamodbav:"description" json:"description"`
	Hint          string              `dynamodbav:"hint" json:"hint"`
	Section       string              `dynamodbav:"section" json:"section"`
	Placeholder   string              `dynamodbav:"placeholder" json:"placeholder"`
	Configs       []map[string]string `dynamodbav:"configs" json:"configs"`
	Options       []Option            `dynamodbav:"options" json:"options"`
	Alert         []Alert             `dynamodbav:"alert" json:"alert"`
	Validations   []map[string]string `dynamodbav:"validations" json:"validations"`
	Order         int                 `dynamodbav:"order" json:"order"`
	Required      bool                `dynamodbav:"required" json:"required"`
	Slug          string              `dynamodbav:"slug" json:"slug"`
	Logic         []FieldLogic        `dynamodbav:"logic" json:"logic"`
	CreatedAt     time.Time           `dynamodbav:"createdAt" json:"createdAt"`
	UpdatedAt     time.Time           `dynamodbav:"updatedAt" json:"updatedAt"`
}

type FieldLogic struct {
	PK                 string      `dynamodbav:"PK"`
	SK                 string      `dynamodbav:"SK"`
	FormVersionFieldID string      `dynamodbav:"formVersionFieldID" json:"formVersionFieldID"` // the field it depends on
	Operator           string      `dynamodbav:"operator" json:"operator"`                     // e.g. equals, not_equals, greater_than
	Value              interface{} `dynamodbav:"value" json:"value"`                           // the value to compare
	Action             string      `dynamodbav:"action" json:"action"`                         // e.g. show, hide, enable, disable
	CreatedAt          time.Time   `dynamodbav:"createdAt" json:"createdAt"`
	UpdatedAt          time.Time   `dynamodbav:"updatedAt" json:"updatedAt"`
}
