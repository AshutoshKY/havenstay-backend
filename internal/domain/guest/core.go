package guest

import (
	"context"

	pb "github.com/user/airbnb-test/api/proto/v1"
)

type Service interface {
	CreateGuest(ctx context.Context, req *pb.CreateGuestRequest) (*pb.GuestResponse, error)
	GetGuest(ctx context.Context, req *pb.GetGuestRequest) (*pb.GuestResponse, error)
}

type coreService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &coreService{repo: repo}
}

func (s *coreService) CreateGuest(ctx context.Context, req *pb.CreateGuestRequest) (*pb.GuestResponse, error) {
	entity := ToEntity(req)

	createdEntity, err := s.repo.Create(ctx, entity)
	if err != nil {
		return nil, err
	}

	return &pb.GuestResponse{
		Guest: EntityToProto(createdEntity),
	}, nil
}

func (s *coreService) GetGuest(ctx context.Context, req *pb.GetGuestRequest) (*pb.GuestResponse, error) {
	entity, err := s.repo.GetByID(ctx, req.GuestId)
	if err != nil {
		return nil, err
	}

	return &pb.GuestResponse{
		Guest: EntityToProto(entity),
	}, nil
}
