package models

import "time"

type Tag struct {
	PK        string    `dynamodbav:"PK"`
	SK        string    `dynamodbav:"SK"`
	CrmID     string    `json:"crmID" dynamodbav:"crmID"`
	CreatedBy string    `json:"createdBy" dynamodbav:"createdBy"`
	Status    string    `json:"status" dynamodbav:"status"`
	Tag       string    `json:"tag" dynamodbav:"tag"`
	CreatedAt time.Time `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" dynamodbav:"updatedAt"`
}
