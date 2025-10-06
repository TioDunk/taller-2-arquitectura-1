package dto

// CreateOrderRequest representa la petición para crear un pedido
type CreateOrderRequest struct {
	CustomerID string `json:"customer_id"`
}

// AddItemRequest representa la petición para agregar un item a un pedido
type AddItemRequest struct {
	ProductID string          `json:"product_id"`
	Quantity  QuantityRequest `json:"quantity"`
}

// OrderItemResponse representa un item del pedido en la respuesta
type OrderItemResponse struct {
	ID        string           `json:"id"`
	Product   ProductResponse  `json:"product"`
	Quantity  QuantityResponse `json:"quantity"`
	UnitPrice MoneyResponse    `json:"unit_price"`
	Total     MoneyResponse    `json:"total"`
}

// OrderResponse representa la respuesta de un pedido
type OrderResponse struct {
	ID         string              `json:"id"`
	CustomerID string              `json:"customer_id"`
	Items      []OrderItemResponse `json:"items"`
	Status     string              `json:"status"`
	Total      MoneyResponse       `json:"total"`
	CreatedAt  string              `json:"created_at"`
	UpdatedAt  string              `json:"updated_at"`
}

// UpdateOrderStatusRequest representa la petición para actualizar el estado de un pedido
type UpdateOrderStatusRequest struct {
	Status string `json:"status"`
}
