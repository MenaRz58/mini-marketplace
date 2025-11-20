package service

import (
	"context"
	"mini-marketplace/products/internal/controller/product"
	"mini-marketplace/products/internal/pkg/model"
	pb "mini-marketplace/proto/products"
)

type Server struct {
	ctrl *product.Controller
	pb.UnimplementedProductsServiceServer
}

func New(ctrl *product.Controller) *Server {
	return &Server{ctrl: ctrl}
}

func (s *Server) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	list := s.ctrl.List()
	resp := &pb.ListProductsResponse{
		Products: make([]*pb.Product, len(list)),
	}
	for i, p := range list {
		resp.Products[i] = convert(p)
	}
	return resp, nil
}

func (s *Server) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	p, err := s.ctrl.Get(req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetProductResponse{Product: convert(p)}, nil
}

func (s *Server) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	p := model.Product{
		ID:    req.Product.Id,
		Name:  req.Product.Name,
		Price: req.Product.Price,
		Stock: int(req.Product.Stock),
	}
	if err := s.ctrl.Create(p); err != nil {
		return nil, err
	}
	return &pb.CreateProductResponse{Product: convert(p)}, nil
}

func (s *Server) ReserveProduct(ctx context.Context, req *pb.ReserveProductRequest) (*pb.ReserveProductResponse, error) {
	p, err := s.ctrl.Reserve(req.Id, int(req.Quantity))
	if err != nil {
		return nil, err
	}
	return &pb.ReserveProductResponse{Product: convert(p)}, nil
}

func convert(p model.Product) *pb.Product {
	return &pb.Product{
		Id:    p.ID,
		Name:  p.Name,
		Price: p.Price,
		Stock: int32(p.Stock),
	}
}
