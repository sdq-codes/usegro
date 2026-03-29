package models

import "time"

type Tag struct {
	ID        string    `bson:"_id" json:"id"`
	CrmID     string    `bson:"crmID" json:"crmID"`
	CreatedBy string    `bson:"createdBy" json:"createdBy"`
	Status    string    `bson:"status" json:"status"`
	Tag       string    `bson:"tag" json:"tag"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}
