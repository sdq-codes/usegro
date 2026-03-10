package models

import "time"

type Form struct {
	PK        string    `dynamodbav:"PK"`
	SK        string    `dynamodbav:"SK"`
	CrmID     string    `json:"crmID" dynamodbav:"crmID"`
	Type      string    `json:"type" dynamodbav:"type"`
	CreatedAt time.Time `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" dynamodbav:"updatedAt"`
}

type FormVersion struct {
	PK                string    `dynamodbav:"PK"`
	SK                string    `dynamodbav:"SK"`
	FormID            string    `json:"formID" dynamodbav:"formID"`
	Title             string    `json:"title" dynamodbav:"title"`
	Description       string    `json:"description" dynamodbav:"description"`
	FormVersionStatus string    `json:"formVersionStatus" dynamodbav:"formVersionStatus"`
	PublishedAt       string    `json:"publishedAt" dynamodbav:"publishedAt"`
	CreatedAt         time.Time `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt" dynamodbav:"updatedAt"`
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
