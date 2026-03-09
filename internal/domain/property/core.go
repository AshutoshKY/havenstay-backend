package property

import (
	"context"

	pb "github.com/user/airbnb-test/api/proto/v1"
)

// Service defines the use cases for the Property domain
type Service interface {
	GetProperty(ctx context.Context, req *pb.GetPropertyRequest) (*pb.PropertyResponse, error)
	ListPropertiesByHost(ctx context.Context, req *pb.ListPropertiesByHostRequest) (*pb.ListPropertiesResponse, error)
	CreateProperty(ctx context.Context, req *pb.CreatePropertyRequest) (*pb.PropertyResponse, error)
	ListProperties(ctx context.Context, req *pb.ListPropertiesRequest) (*pb.ListPropertiesResponse, error)
}

type coreService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &coreService{repo: repo}
}

func (s *coreService) CreateProperty(ctx context.Context, req *pb.CreatePropertyRequest) (*pb.PropertyResponse, error) {
	entity := ToEntity(req)

	createdEntity, err := s.repo.Create(ctx, entity)
	if err != nil {
		return nil, err
	}

	return &pb.PropertyResponse{
		Property: EntityToProto(createdEntity),
	}, nil
}

func (s *coreService) GetProperty(ctx context.Context, req *pb.GetPropertyRequest) (*pb.PropertyResponse, error) {
	entity, err := s.repo.GetByID(ctx, req.PropertyId)
	if err != nil {
		return nil, err
	}

	return &pb.PropertyResponse{
		Property: EntityToProto(entity),
	}, nil
}

func (s *coreService) ListPropertiesByHost(ctx context.Context, req *pb.ListPropertiesByHostRequest) (*pb.ListPropertiesResponse, error) {
	entities, err := s.repo.ListByHostID(ctx, req.HostId, req.Location)
	if err != nil {
		return nil, err
	}

	return &pb.ListPropertiesResponse{
		Properties: EntitiesToProtos(entities),
	}, nil
}

func (s *coreService) ListProperties(ctx context.Context, req *pb.ListPropertiesRequest) (*pb.ListPropertiesResponse, error) {
	entities, err := s.repo.ListAll(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.ListPropertiesResponse{
		Properties: EntitiesToProtos(entities),
	}, nil
}
