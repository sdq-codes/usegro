package models

import "time"

type Form struct {
	ID        string    `bson:"_id" json:"id"`
	CrmID     string    `bson:"crmID" json:"crmID"`
	Type      string    `bson:"type" json:"type"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}

type FormVersion struct {
	ID                string    `bson:"_id" json:"id"`
	FormID            string    `bson:"formID" json:"formID"`
	Title             string    `bson:"title" json:"title"`
	Description       string    `bson:"description" json:"description"`
	FormVersionStatus string    `bson:"formVersionStatus" json:"formVersionStatus"`
	PublishedAt       string    `bson:"publishedAt" json:"publishedAt"`
	CreatedAt         time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt         time.Time `bson:"updatedAt" json:"updatedAt"`
}

type CompleteForm struct {
	Version FormVersion        `json:"version"`
	Fields  []FormVersionField `json:"fields"`
}

type FullForm struct {
	Form    Form               `json:"form"`
	Version FormVersion        `json:"version"`
	Fields  []FormVersionField `json:"fields"`
}
