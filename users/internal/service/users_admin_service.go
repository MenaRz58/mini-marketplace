package service

import (
	"context"

	pb "mini-marketplace/proto/users"
	pba "mini-marketplace/proto/users/admin"
	"mini-marketplace/users/internal/controller/user"
)

type UsersAdminGRPCServer struct {
	ctrl *user.Controller
	pba.UnimplementedUsersAdminServiceServer
}

func NewAdmin(ctrl *user.Controller) *UsersAdminGRPCServer {
	return &UsersAdminGRPCServer{ctrl: ctrl}
}

func (s *UsersAdminGRPCServer) ListUsers(ctx context.Context, req *pba.ListUsersRequest) (*pba.ListUsersResponse, error) {
	usersList := s.ctrl.List()

	resp := &pba.ListUsersResponse{
		Users: make([]*pb.User, len(usersList)),
	}

	for i, u := range usersList {
		resp.Users[i] = &pb.User{
			Id:   u.ID,
			Name: u.Name,
		}
	}
	return resp, nil
}

func (s *UsersAdminGRPCServer) GetUser(ctx context.Context, req *pba.GetUserRequest) (*pba.GetUserResponse, error) {
	u, err := s.ctrl.Get(req.Id)
	if err != nil {
		return nil, err
	}

	return &pba.GetUserResponse{
		User: &pb.User{
			Id:   u.ID,
			Name: u.Name,
		},
	}, nil
}
