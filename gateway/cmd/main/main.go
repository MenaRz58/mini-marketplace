package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	// Importa tus protos generados
	pbGateway "mini-marketplace/proto/gateway" // El proto público
	pbOrders "mini-marketplace/proto/orders"   // El proto interno
	pbProducts "mini-marketplace/proto/products"
	pbUsers "mini-marketplace/proto/users" // El proto interno
)

// 1. Definimos el servidor del Gateway
// Este struct almacena las conexiones a los otros microservicios
type gatewayServer struct {
	pbGateway.UnimplementedGatewayServiceServer // Necesario para compatibilidad

	ordersClient   pbOrders.OrdersServiceClient
	usersClient    pbUsers.UsersServiceClient
	productsClient pbProducts.ProductsServiceClient
}

func main() {
	port := 50051 // Puerto estándar para gRPC (puedes usar el que quieras)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Fallo al escuchar en el puerto %d: %v", port, err)
	}

	// ---------------------------------------------------------
	// 2. CONEXIÓN A MICROSERVICIOS INTERNOS (Service Discovery)
	// En Kubernetes, usamos el nombre del servicio DNS.
	// ---------------------------------------------------------

	// Conectar a Orders Service
	// Nota: Asegúrate que el puerto coincida con el que expone tu servicio Orders (ej. 50052)
	ordersConn, err := grpc.NewClient("orders.mini-marketplace.svc.cluster.local:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("No se pudo conectar a Orders: %v", err)
	}
	defer ordersConn.Close()

	// Conectar a Users Service
	// Nota: Según tus logs anteriores, users escuchaba en 8082
	usersConn, err := grpc.NewClient("users.mini-marketplace.svc.cluster.local:50054", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("No se pudo conectar a Users: %v", err)
	}
	defer usersConn.Close()

	productsConn, err := grpc.NewClient("products.mini-marketplace.svc.cluster.local:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("No se pudo conectar a Products: %v", err)
	}
	defer productsConn.Close()

	// ---------------------------------------------------------
	// 3. INICIALIZAR SERVIDOR GRPC PROPIO
	// ---------------------------------------------------------

	grpcServer := grpc.NewServer()

	// Registramos el Gateway, inyectándole los clientes que acabamos de crear
	pbGateway.RegisterGatewayServiceServer(grpcServer, &gatewayServer{
		ordersClient:   pbOrders.NewOrdersServiceClient(ordersConn),
		usersClient:    pbUsers.NewUsersServiceClient(usersConn),
		productsClient: pbProducts.NewProductsServiceClient(productsConn),
	})

	reflection.Register(grpcServer)

	fmt.Printf("🚀 Gateway gRPC escuchando en el puerto %d\n", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Fallo al servir gRPC: %v", err)
	}
}

// ---------------------------------------------------------
// 4. IMPLEMENTACIÓN DE LA LÓGICA (Handlers)
// ---------------------------------------------------------

func (s *gatewayServer) PlaceOrder(ctx context.Context, req *pbGateway.PlaceOrderRequest) (*pbGateway.PlaceOrderResponse, error) {
	// Validar Usuario (Llamada a Users Service)
	valResp, err := s.usersClient.ValidateUser(ctx, &pbUsers.ValidateUserRequest{UserId: req.UserId})

	if err != nil {
		log.Printf("Error contactando Users Service: %v", err)
		return &pbGateway.PlaceOrderResponse{Success: false}, err
	}

	if !valResp.Valid {
		log.Printf("Usuario inválido o no encontrado: %s", req.UserId)
		// Puedes retornar un error personalizado aquí
		return &pbGateway.PlaceOrderResponse{Success: false}, fmt.Errorf("usuario no valido")
	}

	// PASO B: Preparar datos para Orders
	var orderProducts []*pbOrders.OrderProduct
	for _, item := range req.Items {
		prodInfo, err := s.productsClient.GetProduct(ctx, &pbProducts.GetProductRequest{
			Id: item.ProductId,
		})

		if err != nil {
			log.Printf("Error obteniendo producto %s: %v", item.ProductId, err)
			return &pbGateway.PlaceOrderResponse{Success: false}, fmt.Errorf("producto no encontrado: %s", item.ProductId)
		}

		_, err = s.productsClient.ReserveProduct(ctx, &pbProducts.ReserveProductRequest{
			Id:       item.ProductId,
			Quantity: item.Quantity,
		})

		if err != nil {
			log.Printf("Error reservando stock: %v", err)
			return &pbGateway.PlaceOrderResponse{Success: false}, fmt.Errorf("stock insuficiente o error")
		}

		// B. Usar el precio que nos dio el servicio de productos
		orderProducts = append(orderProducts, &pbOrders.OrderProduct{
			ProductId: item.ProductId,
			Quantity:  item.Quantity,
			Price:     prodInfo.Product.Price,
		})
	}

	// PASO C: Crear la Orden (Llamada a Orders Service)
	createResp, err := s.ordersClient.CreateOrder(ctx, &pbOrders.CreateOrderRequest{
		UserId:   req.UserId,
		Products: orderProducts,
	})

	if err != nil {
		log.Printf("Error creando orden: %v", err)
		return &pbGateway.PlaceOrderResponse{Success: false}, err
	}

	var publicProducts []*pbGateway.OrderProduct
	for _, p := range createResp.Order.Products {
		publicProducts = append(publicProducts, &pbGateway.OrderProduct{
			ProductId: p.ProductId,
			Quantity:  p.Quantity,
			Price:     p.Price,
		})
	}

	publicOrder := &pbGateway.Order{
		Id:        createResp.Order.Id,
		UserId:    createResp.Order.UserId,
		Total:     createResp.Order.Total,
		CreatedAt: createResp.Order.CreatedAt,
		Products:  publicProducts,
	}

	log.Printf("Orden creada exitosamente: %s", createResp.Order.Id)

	// PASO D: Responder al Cliente Externo
	return &pbGateway.PlaceOrderResponse{
		Success: true,
		Order:   publicOrder,
	}, nil
}
