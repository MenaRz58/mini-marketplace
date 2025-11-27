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
		User: &pb.User{
			Id:      u.ID,
			Name:    u.Name,
			Email:   u.Email,
			IsAdmin: u.IsAdmin,
		},
	}, nil
}

func (s *GRPCServer) GetMyProfile(ctx context.Context, req *pb.GetMyProfileRequest) (*pb.GetMyProfileResponse, error) {
	u, err := s.ctrl.Get(req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetMyProfileResponse{
		User: &pb.User{
			Id:      u.ID,
			Name:    u.Name,
			Email:   u.Email,
			IsAdmin: u.IsAdmin,
		},
	}, nil
}

func (s *GRPCServer) ValidateUser(ctx context.Context, req *pb.ValidateUserRequest) (*pb.ValidateUserResponse, error) {
	u, err := s.ctrl.Get(req.UserId)

	if err != nil {

		return &pb.ValidateUserResponse{
			Valid:   false,
			IsAdmin: false,
		}, nil
	}

	return &pb.ValidateUserResponse{
		Valid:   true,
		IsAdmin: u.IsAdmin,
	}, nil
}

func (s *GRPCServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {

	u, err := s.ctrl.Login(ctx, req.UserId, req.Password)

	if err != nil {
		return &pb.LoginResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.LoginResponse{
		Success: true,
		Message: "Login Exitoso",
		User: &pb.User{
			Id:      u.ID,
			Name:    u.Name,
			Email:   u.Email,
			IsAdmin: u.IsAdmin,
		},
	}, nil
}
