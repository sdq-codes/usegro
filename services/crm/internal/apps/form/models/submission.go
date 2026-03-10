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
	PK            string                 `dynamodbav:"PK"`
	SK            string                 `dynamodbav:"SK"`
	FormID        string                 `dynamodbav:"formID"`
	CrmID         string                 `dynamodbav:"crmID"`
	FormVersionID string                 `dynamodbav:"formVersionID"`
	SubmissionID  string                 `dynamodbav:"submissionID"`
	Status        SubmissionStatus       `dynamodbav:"status"`
	Type          SubmissionType         `dynamodbav:"type"`
	Answers       map[string]interface{} `dynamodbav:"answers"`
	VersionSnap   interface{}            `dynamodbav:"versionSnap"`
	CreatedAt     time.Time              `dynamodbav:"createdAt"`
	ArchivedAt    *time.Time             `dynamodbav:"archivedAt,omitempty"`
}
