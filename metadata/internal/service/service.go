package service

import (
	"context"

	"mini-marketplace/metadata/internal/controller"
	"mini-marketplace/metadata/pkg/model"
	pb "mini-marketplace/proto/metadata"
)

type Server struct {
	ctrl *controller.Controller
	pb.UnimplementedMetadataServiceServer
}

func New(ctrl *controller.Controller) *Server {
	return &Server{ctrl: ctrl}
}

// GET
func (s *Server) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	m, err := s.ctrl.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.GetResponse{
		Item: &pb.MetadataItem{
			Id:         m.ID,
			EntityId:   m.EntityID,
			EntityType: m.EntityType,
			Attributes: m.Attributes,
		},
	}, nil
}

// PUT
func (s *Server) Put(ctx context.Context, req *pb.PutRequest) (*pb.PutResponse, error) {

	item := req.Item

	m := &model.Metadata{
		ID:         item.Id,
		EntityID:   item.EntityId,
		EntityType: item.EntityType,
		Attributes: item.Attributes,
	}

	if err := s.ctrl.Put(ctx, m.ID, m); err != nil {
		return nil, err
	}

	return &pb.PutResponse{Success: true}, nil
}
