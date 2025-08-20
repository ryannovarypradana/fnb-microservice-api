package dto

// OrderItemRequest adalah DTO untuk setiap item dalam pesanan.
type OrderItemRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,gt=0"`
}

// CreateOrderRequest adalah DTO untuk membuat pesanan baru.
type CreateOrderRequest struct {
	StoreID uint               `json:"store_id" binding:"required"`
	Items   []OrderItemRequest `json:"items" binding:"required,min=1"`
}
