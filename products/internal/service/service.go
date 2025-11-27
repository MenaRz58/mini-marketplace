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
	productsList, err := s.ctrl.List()
	if err != nil {
		return nil, err
	}
	var responseProducts []*pb.Product
	for _, p := range productsList {
		responseProducts = append(responseProducts, &pb.Product{
			Id:    int32(p.ID),
			Name:  p.Name,
			Price: p.Price,
			Stock: int32(p.Stock),
		})
	}

	return &pb.ListProductsResponse{
		Products: responseProducts,
	}, nil
}

func (s *Server) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	p, err := s.ctrl.Get(uint(req.Id))
	if err != nil {
		return nil, err
	}
	return &pb.GetProductResponse{Product: convert(p)}, nil
}

func (s *Server) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {

	p := model.Product{
		Name:  req.Product.Name,
		Price: req.Product.Price,
		Stock: int(req.Product.Stock),
	}
	if err := s.ctrl.Create(&p); err != nil {
		return nil, err
	}
	return &pb.CreateProductResponse{
		Product: &pb.Product{
			Id:    int32(p.ID),
			Name:  p.Name,
			Price: p.Price,
			Stock: int32(p.Stock),
		},
	}, nil
}

func (s *Server) ReserveProduct(ctx context.Context, req *pb.ReserveProductRequest) (*pb.ReserveProductResponse, error) {
	p, err := s.ctrl.Reserve(uint(req.Id), int(req.Quantity))
	if err != nil {
		return nil, err
	}
	return &pb.ReserveProductResponse{Product: convert(p)}, nil
}

func (s *Server) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error) {
	p := model.Product{
		ID:    uint(req.Product.Id),
		Name:  req.Product.Name,
		Price: req.Product.Price,
		Stock: int(req.Product.Stock),
	}

	if err := s.ctrl.Update(&p); err != nil {
		return nil, err
	}
	return &pb.UpdateProductResponse{
		Product: convert(p),
	}, nil
}

func (s *Server) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	if err := s.ctrl.Delete(uint(req.Id)); err != nil {
		return &pb.DeleteProductResponse{Success: false}, err
	}
	return &pb.DeleteProductResponse{Success: true}, nil
}

func convert(p model.Product) *pb.Product {
	return &pb.Product{
		Id:    int32(p.ID),
		Name:  p.Name,
		Price: p.Price,
		Stock: int32(p.Stock),
	}
}
