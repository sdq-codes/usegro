package models

import "time"

type SubmissionType string

const (
	SubmissionTypeForm            SubmissionType = "form"
	SubmissionTypeCustomer        SubmissionType = "customer"
	SubmissionTypeInvoiceTemplate SubmissionType = "invoice-template"
	SubmissionTypeInvoice         SubmissionType = "invoice"
)

type SubmissionStatus string

const (
	SubmissionStatusActive   SubmissionStatus = "active"
	SubmissionStatusArchived SubmissionStatus = "archived"
)

type FormSubmission struct {
	SubmissionID  string                 `bson:"_id" json:"submissionID"`
	FormID        string                 `bson:"formID" json:"formID"`
	CrmID         string                 `bson:"crmID" json:"crmID"`
	FormVersionID string                 `bson:"formVersionID" json:"formVersionID"`
	Status        SubmissionStatus       `bson:"status" json:"status"`
	Type          SubmissionType         `bson:"type" json:"type"`
	Answers       map[string]interface{} `bson:"answers" json:"answers"`
	VersionSnap   interface{}            `bson:"versionSnap" json:"versionSnap"`
	CreatedAt     time.Time              `bson:"createdAt" json:"createdAt"`
	ArchivedAt    *time.Time             `bson:"archivedAt,omitempty" json:"archivedAt,omitempty"`
}
