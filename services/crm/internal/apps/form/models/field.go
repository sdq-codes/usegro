package models

import "time"

type Option struct {
	Label string `bson:"label" json:"label"`
	Value string `bson:"value" json:"value"`
}

type Alert struct {
	Icon    string `bson:"icon" json:"icon"`
	Type    string `bson:"type" json:"type"`
	Message string `bson:"message" json:"message"`
}

type FormVersionField struct {
	ID            string              `bson:"_id" json:"id"`
	FormVersionID string              `bson:"formVersionID" json:"formVersionID"`
	FieldTypeID   uint                `bson:"fieldTypeId" json:"fieldTypeID"`
	FieldTypeName string              `bson:"fieldTypeName" json:"fieldTypeName"`
	Label         string              `bson:"label" json:"label"`
	Description   string              `bson:"description" json:"description"`
	Hint          string              `bson:"hint" json:"hint"`
	Section       string              `bson:"section" json:"section"`
	Placeholder   string              `bson:"placeholder" json:"placeholder"`
	Configs       []map[string]string `bson:"configs" json:"configs"`
	Options       []Option            `bson:"options" json:"options"`
	Alert         []Alert             `bson:"alert" json:"alert"`
	Validations   []map[string]string `bson:"validations" json:"validations"`
	Order         int                 `bson:"order" json:"order"`
	Required      bool                `bson:"required" json:"required"`
	Slug          string              `bson:"slug" json:"slug"`
	Logic         []FieldLogic        `bson:"logic" json:"logic"`
	CreatedAt     time.Time           `bson:"createdAt" json:"createdAt"`
	UpdatedAt     time.Time           `bson:"updatedAt" json:"updatedAt"`
}

type FieldLogic struct {
	ID                 string      `bson:"_id,omitempty" json:"id,omitempty"`
	FormVersionFieldID string      `bson:"formVersionFieldID" json:"formVersionFieldID"`
	Operator           string      `bson:"operator" json:"operator"`
	Value              interface{} `bson:"value" json:"value"`
	Action             string      `bson:"action" json:"action"`
	CreatedAt          time.Time   `bson:"createdAt" json:"createdAt"`
	UpdatedAt          time.Time   `bson:"updatedAt" json:"updatedAt"`
}
