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
	PK            string                 `dynamodbav:"PK"` // FORM#<formID>
	SK            string                 `dynamodbav:"SK"` // SUBMISSION#<submissionID>
	FormID        string                 `dynamodbav:"formID"`
	CrmID         string                 `dynamodbav:"crmID"`
	FormVersionID string                 `dynamodbav:"formVersionID"`
	SubmissionID  string                 `dynamodbav:"submissionID"`
	Status        SubmissionStatus       `dynamodbav:"status"` // active | archived
	Type          SubmissionType         `dynamodbav:"type"`   // form | customer | invoice-template | invoice
	Answers       map[string]interface{} `dynamodbav:"answers"`
	VersionSnap   interface{}            `dynamodbav:"versionSnap"`
	CreatedAt     time.Time              `dynamodbav:"createdAt"`
	ArchivedAt    *time.Time             `dynamodbav:"archivedAt,omitempty"`
}
