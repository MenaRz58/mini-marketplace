package service

import (
	"context"
	pb "mini-marketplace/proto/users"
	"mini-marketplace/users/internal/controller/user"
)

type GRPCServer struct {
	ctrl *user.Controller
	pb.UnimplementedUsersServiceServer
}

func New(ctrl *user.Controller) *GRPCServer {
	return &GRPCServer{ctrl: ctrl}
}

func (s *GRPCServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	u, err := s.ctrl.Create(req.Id, req.Name)
	if err != nil {
		return nil, err
	}
	return &pb.CreateUserResponse{
		User: &pb.User{Id: u.ID, Name: u.Name},
	}, nil
}

func (s *GRPCServer) GetMyProfile(ctx context.Context, req *pb.GetMyProfileRequest) (*pb.GetMyProfileResponse, error) {
	u, err := s.ctrl.Get(req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetMyProfileResponse{
		User: &pb.User{Id: u.ID, Name: u.Name},
	}, nil
}

// ✅ Implementación de ValidateUser
func (s *GRPCServer) ValidateUser(ctx context.Context, req *pb.ValidateUserRequest) (*pb.ValidateUserResponse, error) {
	isValid, name := s.ctrl.Validate(req.UserId)
	return &pb.ValidateUserResponse{
		Valid: isValid,
		Name:  name,
	}, nil
}
