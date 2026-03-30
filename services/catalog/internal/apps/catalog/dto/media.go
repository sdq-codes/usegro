package dto

import "github.com/google/uuid"

type PresignFileRequest struct {
	Name        string `json:"name"`
	ContentType string `json:"content_type"`
}

type PresignRequest struct {
	Files []PresignFileRequest `json:"files"`
}

type PresignResponse struct {
	Key       string `json:"key"`
	UploadURL string `json:"upload_url"`
	PublicURL string `json:"public_url"`
}

type UploadMediaResponse struct {
	ID       uuid.UUID `json:"id"`
	URL      string    `json:"url"`
	Key      string    `json:"key"`
	MimeType string    `json:"mime_type"`
	Size     int64     `json:"size"`
	Position int       `json:"position"`
}
