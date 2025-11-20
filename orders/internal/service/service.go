package service

import (
	"context"
	"mini-marketplace/orders/internal/controller/order"
	"mini-marketplace/orders/internal/pkg/model"
	pb "mini-marketplace/proto/orders"
)

type Server struct {
	ctrl *order.Controller
	pb.UnimplementedOrdersServiceServer
}

func New(ctrl *order.Controller) *Server {
	return &Server{ctrl: ctrl}
}

func (s *Server) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	products := make([]model.OrderProduct, len(req.Products))
	for i, p := range req.Products {
		products[i] = model.OrderProduct{
			ProductID: p.ProductId,
			Quantity:  int(p.Quantity),
			Price:     p.Price,
		}
	}

	order, err := s.ctrl.Create(req.UserId, products)
	if err != nil {
		return nil, err
	}

	return &pb.CreateOrderResponse{Order: convert(order)}, nil
}

func (s *Server) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	order, err := s.ctrl.Get(req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetOrderResponse{Order: convert(order)}, nil
}

func (s *Server) ListOrders(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	list, err := s.ctrl.List()
	if err != nil {
		return nil, err
	}

	resp := &pb.ListResponse{Orders: make([]*pb.Order, len(list))}
	for i, o := range list {
		resp.Orders[i] = convert(o)
	}
	return resp, nil
}

// convert convierte model.Order a pb.Order
func convert(o model.Order) *pb.Order {
	prods := make([]*pb.OrderProduct, len(o.Products))
	for i, p := range o.Products {
		prods[i] = &pb.OrderProduct{
			ProductId: p.ProductID,
			Quantity:  int32(p.Quantity),
			Price:     p.Price,
		}
	}
	return &pb.Order{
		Id:        o.ID,
		UserId:    o.UserID,
		Products:  prods,
		Total:     o.Total,
		CreatedAt: o.CreatedAt,
	}
}
