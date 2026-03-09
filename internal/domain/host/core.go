package host

import (
	"context"

	pb "github.com/user/airbnb-test/api/proto/v1"
)

// Service defines the use cases for the Host domain
type Service interface {
	CreateHost(ctx context.Context, req *pb.CreateHostRequest) (*pb.HostResponse, error)
	GetHost(ctx context.Context, req *pb.GetHostRequest) (*pb.HostResponse, error)
}

type coreService struct {
	repo Repository
}

// NewService creates a new Host core service
func NewService(repo Repository) Service {
	return &coreService{repo: repo}
}

func (s *coreService) CreateHost(ctx context.Context, req *pb.CreateHostRequest) (*pb.HostResponse, error) {
	// 1. Validate request (skipped for brevity)

	// 2. Map to domain entity
	entity := ToEntity(req)

	// Default to unverified on creation
	entity.IsVerified = false

	// 3. Save via repo
	createdEntity, err := s.repo.Create(ctx, entity)
	if err != nil {
		return nil, err
	}

	// 4. Map back to proto
	return &pb.HostResponse{
		Host: EntityToProto(createdEntity),
	}, nil
}

func (s *coreService) GetHost(ctx context.Context, req *pb.GetHostRequest) (*pb.HostResponse, error) {
	entity, err := s.repo.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.HostResponse{
		Host: EntityToProto(entity),
	}, nil
}
