package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sdq-codes/usegro-api/config"
	"github.com/sdq-codes/usegro-api/internal/apps/base/models"
	"github.com/sdq-codes/usegro-api/internal/apps/base/repositories"
	"github.com/sdq-codes/usegro-api/internal/helper/auth"
	grpcclient "github.com/sdq-codes/usegro-api/internal/helper/grpc"
	"github.com/sdq-codes/usegro-api/pkg/exception"
	pb "github.com/usegro/proto/crm"
	"gorm.io/gorm"
)

type Service struct {
	db             *gorm.DB
	userRepository *repositories.UserRepository
}

func NewUserService(db *gorm.DB) *Service {
	return &Service{
		db:             db,
		userRepository: repositories.NewUserRepository(),
	}
}

func (s *Service) GetLoggedInUser(fCtx *fiber.Ctx) (*models.User, []*pb.CrmOrganization, error) {
	claims, err := auth.AuthUser(fCtx)
	if err != nil {
		return nil, nil, exception.UnauthorizedError
	}

	userID := claims.User.ID.String()

	user, err := s.userRepository.GetUserById(fCtx.UserContext(), s.db, userID)
	if err != nil {
		return nil, nil, exception.UnauthorizedError
	}

	crmAddr := config.GetConfig().CrmService.Address
	crmClient, err := grpcclient.NewCRMClient(crmAddr)
	var crms []*pb.CrmOrganization
	if err == nil {
		crms, _ = crmClient.ListOrganizationsByUser(fCtx.UserContext(), userID)
	}

	return user, crms, nil
}
