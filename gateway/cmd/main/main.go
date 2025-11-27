package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pbOrders "mini-marketplace/proto/orders"
	pbProducts "mini-marketplace/proto/products"
	pbUsers "mini-marketplace/proto/users"
)

type BuyRequest struct {
	UserID string `json:"user_id"`
	Items  []struct {
		ProductID int32 `json:"product_id"`
		Quantity  int   `json:"quantity"`
	} `json:"items"`
}

type LoginRequest struct {
	UserID   string `json:"user_id"`
	Password string `json:"password"`
}

var (
	ordersClient   pbOrders.OrdersServiceClient
	usersClient    pbUsers.UsersServiceClient
	productsClient pbProducts.ProductsServiceClient
)

func main() {
	port := ":9090"

	// Conexiones gRPC
	connP, _ := grpc.NewClient("products.mini-marketplace.svc.cluster.local:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	productsClient = pbProducts.NewProductsServiceClient(connP)
	connO, _ := grpc.NewClient("orders.mini-marketplace.svc.cluster.local:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	ordersClient = pbOrders.NewOrdersServiceClient(connO)
	connU, _ := grpc.NewClient("users.mini-marketplace.svc.cluster.local:50054", grpc.WithTransportCredentials(insecure.NewCredentials()))
	usersClient = pbUsers.NewUsersServiceClient(connU)

	// 1. Assets
	assetsFs := http.FileServer(http.Dir("/app/static"))
	http.Handle("/assets/", assetsFs)
	http.Handle("/vite.svg", assetsFs)

	// 2. Handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := "/app/static/index.html"

		log.Printf("Intentando leer archivo: %s", path)

		content, err := os.ReadFile(path)
		if err != nil {
			log.Printf("ERROR LEYENDO ARCHIVO: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Println("✅ Archivo leído con éxito. Enviando al navegador.")
		w.Header().Set("Content-Type", "text/html")
		w.Write(content)
	})

	// 3. API
	http.HandleFunc("/api/login", handleLogin)
	http.HandleFunc("/api/buy", handleBuy)
	http.HandleFunc("/api/products", handleListProducts)
	http.HandleFunc("/api/orders", handleListOrders)
	http.HandleFunc("/api/admin/create-product", handleCreateProduct)
	http.HandleFunc("/api/admin/update-product", handleUpdateProduct)
	http.HandleFunc("/api/admin/delete-product", handleDeleteProduct)

	log.Printf("SERVIDOR EN FRECUENCIA NUEVA - PUERTO %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", 400)
		return
	}

	ctx := context.Background()
	loginResp, err := usersClient.Login(ctx, &pbUsers.LoginRequest{
		UserId:   req.UserID,
		Password: req.Password,
	})

	if err != nil {
		log.Printf("Error gRPC Login: %v", err)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "error": "Servicio de usuarios no disponible"})
		return
	}

	if !loginResp.Success {
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "error": loginResp.Message})
		return
	}

	log.Printf("Usuario autenticado correctamente: %s (%s)", loginResp.User.Name, req.UserID)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"user": map[string]interface{}{
			"id":       loginResp.User.Id,
			"name":     loginResp.User.Name,
			"email":    loginResp.User.Email,
			"is_admin": loginResp.User.IsAdmin,
		},
	})
}

func handleListProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp, err := productsClient.ListProducts(context.Background(), &pbProducts.ListProductsRequest{})
	if err != nil {
		http.Error(w, "Error interno", 500)
		return
	}
	json.NewEncoder(w).Encode(resp.Products)
}

func handleBuy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	var req BuyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", 400)
		return
	}
	ctx := context.Background()
	valResp, err := usersClient.ValidateUser(ctx, &pbUsers.ValidateUserRequest{UserId: req.UserID})
	if err != nil || !valResp.Valid {
		jsonResponse(w, false, nil, "Usuario no válido")
		return
	}
	var orderProducts []*pbOrders.OrderProduct
	for _, item := range req.Items {
		resResp, err := productsClient.ReserveProduct(ctx, &pbProducts.ReserveProductRequest{
			Id:       item.ProductID,
			Quantity: int32(item.Quantity),
		})
		if err != nil {
			jsonResponse(w, false, nil, "Stock insuficiente")
			return
		}
		orderProducts = append(orderProducts, &pbOrders.OrderProduct{
			ProductId: item.ProductID,
			Quantity:  int32(item.Quantity),
			Price:     resResp.Product.Price,
		})
	}
	createResp, err := ordersClient.CreateOrder(ctx, &pbOrders.CreateOrderRequest{
		UserId:   req.UserID,
		Products: orderProducts,
	})
	if err != nil {
		jsonResponse(w, false, nil, "Error al crear orden")
		return
	}
	jsonResponse(w, true, createResp.Order, "")
}

func handleListOrders(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "Falta user_id", 400)
		return
	}

	resp, err := ordersClient.ListOrders(context.Background(), &pbOrders.ListRequest{})
	if err != nil {
		log.Printf("Error trayendo órdenes: %v", err)
		http.Error(w, "Error obteniendo órdenes", 500)
		return
	}

	var userOrders []map[string]interface{}

	for _, o := range resp.Orders {
		if o.UserId == userID {

			var items []map[string]interface{}
			for _, p := range o.Products {
				items = append(items, map[string]interface{}{
					"id":       p.ProductId,
					"quantity": p.Quantity,
					"price":    p.Price,
				})
			}

			userOrders = append(userOrders, map[string]interface{}{
				"id":    o.Id,
				"total": o.Total,
				"date":  o.CreatedAt,
				"items": items,
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userOrders)
}

func handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		UserID string  `json:"user_id"`
		Name   string  `json:"name"`
		Price  float32 `json:"price"`
		Stock  int32   `json:"stock"`
	}

	valResp, err := usersClient.ValidateUser(context.Background(), &pbUsers.ValidateUserRequest{UserId: req.UserID})

	if err != nil || !valResp.Valid {
		http.Error(w, "Usuario inválido", http.StatusForbidden)
		return
	}

	if !valResp.IsAdmin {
		http.Error(w, "ACCESO DENEGADO: Se requieren permisos de Administrador", http.StatusForbidden)
		return
	}

	resp, err := productsClient.CreateProduct(context.Background(), &pbProducts.CreateProductRequest{
		Product: &pbProducts.Product{
			Name:  req.Name,
			Price: float64(req.Price),
			Stock: int32(req.Stock),
		},
	})

	if err != nil {
		http.Error(w, "Error creando producto: "+err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"product": resp.Product,
	})
}

func handleUpdateProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		UserID string  `json:"user_id"`
		ID     int32   `json:"id"`
		Name   string  `json:"name"`
		Price  float32 `json:"price"`
		Stock  int32   `json:"stock"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	if req.UserID != "u1" {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	resp, err := productsClient.UpdateProduct(context.Background(), &pbProducts.UpdateProductRequest{
		Product: &pbProducts.Product{
			Id:    req.ID,
			Name:  req.Name,
			Price: float64(req.Price),
			Stock: req.Stock,
		},
	})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"product": resp.Product,
	})
}

func handleDeleteProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		UserID string `json:"user_id"`
		ID     int32  `json:"id"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	if req.UserID != "u1" {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	_, err := productsClient.DeleteProduct(context.Background(), &pbProducts.DeleteProductRequest{Id: req.ID})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func jsonResponse(w http.ResponseWriter, success bool, order *pbOrders.Order, errMsg string) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"success": success,
	}
	if order != nil {
		response["order"] = map[string]interface{}{
			"id":    order.Id,
			"total": order.Total,
		}
	}
	if errMsg != "" {
		response["error"] = errMsg
	}
	json.NewEncoder(w).Encode(response)
}
