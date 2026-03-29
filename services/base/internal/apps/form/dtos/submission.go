package dtos

type CreateSubmissionInput struct {
	Answers     map[string]interface{} `json:"answers" validate:"required"`
	VersionSnap interface{}            `json:"versionSnap" validate:"required"`
	Type        string                 `json:"type" validate:"required"`
}

type UpdateSubmissionStatusInput struct {
	Status string `json:"status" validate:"required,oneof=active archived"`
}

type UpdateSubmissionAnswersInput struct {
	Answers map[string]interface{} `json:"answers" validate:"required"`
}
