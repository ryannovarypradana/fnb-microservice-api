// internal/order/service.go
package order

import (
	"encoding/json"
	"errors"
	"fmt"
	"fnb-system/pkg/dto"
	"fnb-system/pkg/eventbus"
	"fnb-system/pkg/model"
	"net/http"
)

// ProductClient adalah interface untuk berkomunikasi dengan Product Service.
type ProductClient interface {
	GetMenuByID(id uint) (*model.Menu, error)
}

// OrderService mendefinisikan kontrak untuk logika bisnis pesanan.
type OrderService interface {
	CreateOrder(userID uint, req dto.CreateOrderRequest) (*model.Order, error)
	GetMyOrders(userID uint) (*[]model.Order, error)
}

// orderService adalah implementasi dari OrderService.
type orderService struct {
	repo          OrderRepository
	productClient ProductClient
	eventBus      eventbus.EventBus // Dependensi untuk RabbitMQ
}

// NewOrderService membuat instance baru dari orderService.
func NewOrderService(repo OrderRepository, productClient ProductClient, bus eventbus.EventBus) OrderService {
	return &orderService{
		repo:          repo,
		productClient: productClient,
		eventBus:      bus,
	}
}

// GetMyOrders mengambil riwayat pesanan milik seorang user.
func (s *orderService) GetMyOrders(userID uint) (*[]model.Order, error) {
	return s.repo.FindOrdersByUserID(userID)
}

// CreateOrder adalah logika inti untuk membuat pesanan.
func (s *orderService) CreateOrder(userID uint, req dto.CreateOrderRequest) (*model.Order, error) {
	var orderItems []model.OrderItem
	var totalPrice float64

	if len(req.Items) == 0 {
		return nil, errors.New("order must have at least one item")
	}

	// 1. Validasi setiap item dan hitung total harga
	for _, itemReq := range req.Items {
		// Panggil Product Service untuk mendapatkan detail menu
		menu, err := s.productClient.GetMenuByID(itemReq.MenuID)
		if err != nil {
			return nil, fmt.Errorf("menu with id %d not found or product service is down", itemReq.MenuID)
		}

		itemPrice := menu.Price * float64(itemReq.Quantity)
		totalPrice += itemPrice

		orderItems = append(orderItems, model.OrderItem{
			MenuID:   itemReq.MenuID,
			Quantity: itemReq.Quantity,
			Price:    itemPrice, // Simpan harga total per item saat itu
		})
	}

	// 2. Buat objek Order utama
	newOrder := model.Order{
		UserID:     userID,
		TotalPrice: totalPrice,
		Status:     "pending",
	}

	// 3. Simpan ke database menggunakan transaksi
	createdOrder, err := s.repo.CreateOrderInTx(&newOrder, &orderItems)
	if err != nil {
		return nil, err
	}

	// 4. Publikasikan event ke RabbitMQ setelah pesanan berhasil dibuat
	eventPayload := map[string]interface{}{
		"order_id":    createdOrder.ID,
		"user_id":     createdOrder.UserID,
		"total_price": createdOrder.TotalPrice,
		"status":      createdOrder.Status,
	}
	// Menggunakan goroutine agar tidak memblokir response ke user
	go s.eventBus.Publish("orders", "order.created", eventPayload)

	return createdOrder, nil
}

// --- Implementasi HTTP Client untuk Product Service ---

type productClient struct {
	baseURL string
}

func NewProductClient(baseURL string) ProductClient {
	return &productClient{baseURL}
}

func (c *productClient) GetMenuByID(id uint) (*model.Menu, error) {
	// Di Product Service, kita perlu membuat endpoint ini: GET /api/v1/menus/:id
	url := fmt.Sprintf("%s/api/v1/menus/%d", c.baseURL, id)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call product service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("product service returned status %d", resp.StatusCode)
	}

	var menu model.Menu
	if err := json.NewDecoder(resp.Body).Decode(&menu); err != nil {
		return nil, fmt.Errorf("failed to decode product response: %w", err)
	}

	return &menu, nil
}
