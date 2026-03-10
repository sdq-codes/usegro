package dto

import (
	"github.com/google/uuid"
	"github.com/sdq-codes/usegro-api/internal/apps/base/models"
	pb "github.com/usegro/proto/crm"
)

type UserDTO struct {
	ID            uuid.UUID             `json:"id"`
	Email         string                `json:"email"`
	Verifications []models.Verification `json:"verifications"`
	Organizations []*pb.CrmOrganization `json:"organizations"`
}
