package entities

import (
	"ddd-poc/domain/valueobjects"
	"errors"

	"github.com/google/uuid"
)

// OrderItem representa un item dentro de un pedido
// Es parte del Aggregate Order
type OrderItem struct {
	id        string
	product   Product
	quantity  valueobjects.Quantity
	unitPrice valueobjects.Money
}

// NewOrderItem crea un nuevo item de pedido
func NewOrderItem(product Product, quantity valueobjects.Quantity) (*OrderItem, error) {
	if !product.IsActive() {
		return nil, errors.New("cannot add inactive product to order")
	}

	if quantity.Value() <= 0 {
		return nil, errors.New("quantity must be greater than zero")
	}

	return &OrderItem{
		id:        uuid.New().String(),
		product:   product,
		quantity:  quantity,
		unitPrice: product.Price(),
	}, nil
}

// OrderItemFromExisting crea un item de pedido a partir de datos existentes
func OrderItemFromExisting(id string, product Product, quantity valueobjects.Quantity, unitPrice valueobjects.Money) *OrderItem {
	return &OrderItem{
		id:        id,
		product:   product,
		quantity:  quantity,
		unitPrice: unitPrice,
	}
}

// ID retorna el ID único del item
func (oi OrderItem) ID() string {
	return oi.id
}

// Product retorna el producto del item
func (oi OrderItem) Product() Product {
	return oi.product
}

// Quantity retorna la cantidad del item
func (oi OrderItem) Quantity() valueobjects.Quantity {
	return oi.quantity
}

// UnitPrice retorna el precio unitario del item
func (oi OrderItem) UnitPrice() valueobjects.Money {
	return oi.unitPrice
}

// TotalPrice calcula el precio total del item
func (oi OrderItem) TotalPrice() (valueobjects.Money, error) {
	factor := float64(oi.quantity.Value())
	return oi.unitPrice.Multiply(factor)
}

// UpdateQuantity actualiza la cantidad del item
func (oi *OrderItem) UpdateQuantity(newQuantity valueobjects.Quantity) error {
	if newQuantity.Value() <= 0 {
		return errors.New("quantity must be greater than zero")
	}
	oi.quantity = newQuantity
	return nil
}

// Equals compara dos items por su ID
func (oi OrderItem) Equals(other OrderItem) bool {
	return oi.id == other.id
}
