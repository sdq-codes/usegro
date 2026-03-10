package grpcserver

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	pb "github.com/usegro/proto/crm"
	"github.com/usegro/services/crm/internal/apps/crm/repositories"
)

type CRMServiceServer struct {
	pb.UnimplementedCRMServiceServer
	db   *gorm.DB
	repo *repositories.CRMUserOrganizationRepository
}

func NewCRMServiceServer(db *gorm.DB, rdb *redis.Client) *CRMServiceServer {
	repo := repositories.NewCRMUserOrganizationRepository(db, rdb)
	return &CRMServiceServer{db: db, repo: repo}
}

func (s *CRMServiceServer) ListOrganizationsByUser(ctx context.Context, req *pb.ListOrganizationsByUserRequest) (*pb.ListOrganizationsByUserResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	orgs, err := s.repo.FetchCRMUserOrganization(ctx, s.db, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch organizations: %v", err)
	}

	var pbOrgs []*pb.CrmOrganization
	for _, org := range *orgs {
		pbOrgs = append(pbOrgs, &pb.CrmOrganization{
			Id:           org.ID.String(),
			UserId:       org.UserID.String(),
			FullName:     org.FullName,
			BusinessName: org.BusinessName,
			BusinessInfo: string(org.BusinessInfo),
			Active:       org.Active,
			CreatedAt:    timestamppb.New(org.CreatedAt),
		})
	}

	return &pb.ListOrganizationsByUserResponse{Organizations: pbOrgs}, nil
}

func (s *CRMServiceServer) CreateTenant(ctx context.Context, req *pb.CreateTenantRequest) (*pb.Tenant, error) {
	if req.Name == "" || req.WabaId == "" {
		return nil, status.Error(codes.InvalidArgument, "name and waba_id are required")
	}
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (s *CRMServiceServer) GetTenant(ctx context.Context, req *pb.GetTenantRequest) (*pb.Tenant, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (s *CRMServiceServer) UpdateTenant(ctx context.Context, req *pb.UpdateTenantRequest) (*pb.Tenant, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (s *CRMServiceServer) CreateContact(ctx context.Context, req *pb.CreateContactRequest) (*pb.Contact, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (s *CRMServiceServer) GetContact(ctx context.Context, req *pb.GetContactRequest) (*pb.Contact, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (s *CRMServiceServer) UpdateContact(ctx context.Context, req *pb.UpdateContactRequest) (*pb.Contact, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (s *CRMServiceServer) ListContacts(ctx context.Context, req *pb.ListContactsRequest) (*pb.ListContactsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (s *CRMServiceServer) UpsertContactByPhone(ctx context.Context, req *pb.UpsertContactByPhoneRequest) (*pb.Contact, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (s *CRMServiceServer) CreateLead(ctx context.Context, req *pb.CreateLeadRequest) (*pb.Lead, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (s *CRMServiceServer) UpdateLead(ctx context.Context, req *pb.UpdateLeadRequest) (*pb.Lead, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (s *CRMServiceServer) ListLeads(ctx context.Context, req *pb.ListLeadsRequest) (*pb.ListLeadsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (s *CRMServiceServer) UpdateLifecycleStage(ctx context.Context, req *pb.UpdateLifecycleStageRequest) (*pb.Lead, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (s *CRMServiceServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (s *CRMServiceServer) AuthenticateUser(ctx context.Context, req *pb.AuthenticateUserRequest) (*pb.AuthResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (s *CRMServiceServer) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func StartGRPCServer(port int, db *gorm.DB, rdb *redis.Client) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("gRPC listen failed: %v", err)
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(loggingInterceptor))
	srv := NewCRMServiceServer(db, rdb)
	pb.RegisterCRMServiceServer(s, srv)

	log.Printf("CRM gRPC server listening on :%d", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("gRPC serve failed: %v", err)
	}
}

func loggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	log.Printf("[CRM gRPC] %s | %v | err=%v", info.FullMethod, time.Since(start), err)
	return resp, err
}
