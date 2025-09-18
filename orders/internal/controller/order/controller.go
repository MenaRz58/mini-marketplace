package order

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	metadataModel "mini-marketplace/metadata/pkg/model"
	"mini-marketplace/orders/internal/pkg/model"
	"mini-marketplace/pkg"
	"mini-marketplace/pkg/discovery/consul"

	"github.com/google/uuid"
)

type Repo interface {
	List() []model.Order
	Get(id string) (model.Order, error)
	Create(o model.Order) error
}

type Controller struct {
	repo     Repo
	registry *consul.Registry
	ctx      context.Context
}

func NewController(r Repo, reg *consul.Registry, ctx context.Context) *Controller {
	return &Controller{
		repo:     r,
		registry: reg,
		ctx:      ctx,
	}
}

func (c *Controller) List() []model.Order {
	return c.repo.List()
}

func (c *Controller) Get(id string) (model.Order, error) {
	return c.repo.Get(id)
}

// Create soporta múltiples productos por orden
func (c *Controller) Create(o model.Order) error {
	if o.UserID == "" || len(o.Products) == 0 {
		return errors.New("invalid order fields")
	}

	// Generar ID automáticamente si no se proporciona
	if o.ID == "" {
		o.ID = uuid.New().String()
	} else {
		_, err := c.repo.Get(o.ID)
		if err == nil {
			return errors.New("order already exists")
		}
	}

	// Descubrir servicios
	usersAddrs, err := c.registry.ServiceAddress(c.ctx, "users")
	if err != nil || len(usersAddrs) == 0 {
		return errors.New("users service unavailable")
	}
	productsAddrs, err := c.registry.ServiceAddress(c.ctx, "products")
	if err != nil || len(productsAddrs) == 0 {
		return errors.New("products service unavailable")
	}
	metadataAddrs, err := c.registry.ServiceAddress(c.ctx, "metadata")
	if err != nil || len(metadataAddrs) == 0 {
		return errors.New("metadata service unavailable")
	}

	userURL := "http://" + usersAddrs[0]
	productURL := "http://" + productsAddrs[0]
	metadataURL := "http://" + metadataAddrs[0]

	// Validar usuario
	usr, err := fetchUser(userURL, o.UserID)
	if err != nil || usr.ID == "" {
		return errors.New("user not found")
	}

	// Calcular total y reservar stock de cada producto
	total := 0.0
	for i, p := range o.Products {
		prod, err := fetchProduct(productURL, p.ProductID)
		if err != nil || prod.ID == "" {
			return fmt.Errorf("product %s not found", p.ProductID)
		}

		if err := reserveStock(productURL, prod.ID, p.Quantity); err != nil {
			return err
		}

		o.Products[i].Price = prod.Price
		total += prod.Price * float64(p.Quantity)
	}

	o.Total = total
	o.CreatedAt = time.Now().Unix()

	// Guardar orden
	if err := c.repo.Create(o); err != nil {
		return err
	}

	// Registrar metadata
	md := &metadataModel.Metadata{
		ID:         fmt.Sprintf("md-%s", o.ID),
		EntityID:   o.ID,
		EntityType: "order",
		Attributes: map[string]string{
			"user_id":  o.UserID,
			"total":    fmt.Sprintf("%.2f", o.Total),
			"products": fmt.Sprintf("%v", o.Products),
		},
	}
	data, _ := json.Marshal(md)
	_, _ = http.Post(metadataURL+"/metadata", "application/json", bytes.NewReader(data))

	return nil
}

// fetchUser llama al endpoint /users/{id}
func fetchUser(base, id string) (pkg.UserRef, error) {
	resp, err := http.Get(fmt.Sprintf("%s/users/%s", base, id))
	if err != nil {
		return pkg.UserRef{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return pkg.UserRef{}, errors.New("user not ok")
	}
	var u pkg.UserRef
	return u, json.NewDecoder(resp.Body).Decode(&u)
}

// fetchProduct llama al endpoint /products/{id}
func fetchProduct(base, id string) (pkg.ProductRef, error) {
	resp, err := http.Get(fmt.Sprintf("%s/products/%s", base, id))
	if err != nil {
		return pkg.ProductRef{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return pkg.ProductRef{}, errors.New("product not ok")
	}
	var p pkg.ProductRef
	return p, json.NewDecoder(resp.Body).Decode(&p)
}

// reserveStock llama a POST /products/reserve
func reserveStock(base, id string, qty int) error {
	reqBody := struct {
		ID       string `json:"id"`
		Quantity int    `json:"quantity"`
	}{ID: id, Quantity: qty}

	b, _ := json.Marshal(reqBody)
	resp, err := http.Post(fmt.Sprintf("%s/products/reserve", base), "application/json", bytesReader(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		return errors.New("cannot reserve stock")
	}
	return nil
}

// bytesReader crea un *bytes.Reader compatible
func bytesReader(b []byte) *bytesReaderWrapper { return &bytesReaderWrapper{b: b} }

type bytesReaderWrapper struct{ b []byte }

func (r *bytesReaderWrapper) Read(p []byte) (int, error) {
	if len(r.b) == 0 {
		return 0, io.EOF
	}
	n := copy(p, r.b)
	r.b = r.b[n:]
	return n, nil
}
func (r *bytesReaderWrapper) Close() error { return nil }
