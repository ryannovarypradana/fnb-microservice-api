// pkg/dto/order.go
package dto

// OrderItemRequest adalah struct untuk satu item dalam pesanan.
type OrderItemRequest struct {
	MenuID   uint `json:"menu_id"`
	Quantity int  `json:"quantity"`
}

// CreateOrderRequest adalah request body utama untuk membuat pesanan baru.
type CreateOrderRequest struct {
	Items []OrderItemRequest `json:"items"`
}
