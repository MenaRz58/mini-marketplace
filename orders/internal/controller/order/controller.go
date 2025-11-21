package order

import (
	"context"
	"errors"
	"time"

	"mini-marketplace/orders/internal/pkg/model"
	pb "mini-marketplace/proto/orders"

	"github.com/google/uuid"
)

type Repository interface {
	Save(o *model.Order) error
	Get(id string) (*model.Order, error)
	List() ([]model.Order, error)
}

type Controller struct {
	repo Repository
}

func NewController(r Repository) *Controller {
	return &Controller{repo: r}
}

// List devuelve todas las órdenes
func (c *Controller) List() ([]model.Order, error) {
	return c.repo.List()
}

// Get devuelve una orden por ID
func (c *Controller) Get(id string) (*model.Order, error) {
	return c.repo.Get(id)
}

// Create crea una nueva orden
func (c *Controller) Create(userID string, products []model.OrderProduct) (*model.Order, error) {
	if userID == "" || len(products) == 0 {
		return nil, errors.New("invalid order fields")
	}

	o := model.Order{
		ID:        uuid.New().String(),
		UserID:    userID,
		Products:  products,
		CreatedAt: time.Now().Unix(),
	}

	total := 0.0
	for _, p := range products {
		total += p.Price * float64(p.Quantity)
	}
	o.Total = total

	if err := c.repo.Save(&o); err != nil {
		return nil, err
	}

	return &o, nil
}

// --------------------
// gRPC Server Adapter
// --------------------

type GRPCServer struct {
	ctrl *Controller
	pb.UnimplementedOrdersServiceServer
}

func NewGRPCServer(ctrl *Controller) *GRPCServer {
	return &GRPCServer{ctrl: ctrl}
}

// ListOrders implementa el RPC ListOrders
func (s *GRPCServer) ListOrders(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	ordersList, err := s.ctrl.List()
	if err != nil {
		return nil, err
	}

	resp := &pb.ListResponse{}
	for _, o := range ordersList {
		orderProto := &pb.Order{
			Id:        o.ID,
			UserId:    o.UserID,
			Total:     o.Total,
			CreatedAt: o.CreatedAt,
		}
		for _, p := range o.Products {
			orderProto.Products = append(orderProto.Products, &pb.OrderProduct{
				ProductId: p.ProductID,
				Quantity:  int32(p.Quantity),
				Price:     p.Price,
			})
		}
		resp.Orders = append(resp.Orders, orderProto)
	}
	return resp, nil
}

// GetOrder implementa el RPC GetOrder
func (s *GRPCServer) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	o, err := s.ctrl.Get(req.Id)
	if err != nil {
		return nil, err
	}

	orderProto := &pb.Order{
		Id:        o.ID,
		UserId:    o.UserID,
		Total:     o.Total,
		CreatedAt: o.CreatedAt,
	}
	for _, p := range o.Products {
		orderProto.Products = append(orderProto.Products, &pb.OrderProduct{
			ProductId: p.ProductID,
			Quantity:  int32(p.Quantity),
			Price:     p.Price,
		})
	}

	return &pb.GetOrderResponse{Order: orderProto}, nil
}

// CreateOrder implementa el RPC CreateOrder
func (s *GRPCServer) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	o := model.Order{
		UserID: req.UserId,
	}

	for _, p := range req.Products {
		o.Products = append(o.Products, model.OrderProduct{
			ProductID: p.ProductId,
			Quantity:  int32(p.Quantity),
			Price:     p.Price,
		})
	}

	createdOrder, err := s.ctrl.Create(o.UserID, o.Products)
	if err != nil {
		return nil, err
	}

	orderProto := &pb.Order{
		Id:        createdOrder.ID,
		UserId:    createdOrder.UserID,
		Total:     createdOrder.Total,
		CreatedAt: createdOrder.CreatedAt,
	}
	for _, p := range createdOrder.Products {
		orderProto.Products = append(orderProto.Products, &pb.OrderProduct{
			ProductId: p.ProductID,
			Quantity:  int32(p.Quantity),
			Price:     p.Price,
		})
	}

	return &pb.CreateOrderResponse{Order: orderProto}, nil
}
