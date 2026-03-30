package models

import "time"

const ActivityTypeCustomerCreated = "customer_created"

type CustomerActivity struct {
	ID           string    `bson:"_id,omitempty" json:"id,omitempty"`
	ActivityType string    `bson:"activityType" json:"activityType"`
	Description  string    `bson:"description" json:"description"`
	CrmID        string    `bson:"crmID" json:"crmID"`
	CustomerID   string    `bson:"customerID" json:"customerID"`
	FormID       string    `bson:"formID" json:"formID"`
	PerformedBy  string    `bson:"performedBy" json:"performedBy"`
	CreatedAt    time.Time `bson:"createdAt" json:"createdAt"`
}
