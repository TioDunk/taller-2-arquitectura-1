package entities

import (
	"ddd-poc/domain/valueobjects"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// OrderStatus representa el estado de un pedido
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusConfirmed OrderStatus = "confirmed"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusDelivered OrderStatus = "delivered"
	OrderStatusCancelled OrderStatus = "cancelled"
)

// Order representa un pedido en el sistema
// Es un Aggregate Root porque:
// 1. Es la única forma de acceder a los OrderItems
// 2. Mantiene la consistencia del dominio
// 3. Coordina las operaciones sobre sus entidades hijas
type Order struct {
	id         string
	customerID string
	items      []OrderItem
	status     OrderStatus
	createdAt  time.Time
	updatedAt  time.Time
}

// NewOrder crea un nuevo pedido
func NewOrder(customerID string) (*Order, error) {
	if customerID == "" {
		return nil, errors.New("customer ID cannot be empty")
	}

	now := time.Now()
	return &Order{
		id:         uuid.New().String(),
		customerID: customerID,
		items:      []OrderItem{},
		status:     OrderStatusPending,
		createdAt:  now,
		updatedAt:  now,
	}, nil
}

// OrderFromExisting crea un pedido a partir de datos existentes
func OrderFromExisting(id, customerID string, items []OrderItem, status OrderStatus, createdAt, updatedAt time.Time) *Order {
	return &Order{
		id:         id,
		customerID: customerID,
		items:      items,
		status:     status,
		createdAt:  createdAt,
		updatedAt:  updatedAt,
	}
}

// ID retorna el ID único del pedido
func (o Order) ID() string {
	return o.id
}

// CustomerID retorna el ID del cliente
func (o Order) CustomerID() string {
	return o.customerID
}

// Items retorna una copia de los items del pedido
func (o Order) Items() []OrderItem {
	// Retornar una copia para mantener la encapsulación
	items := make([]OrderItem, len(o.items))
	copy(items, o.items)
	return items
}

// Status retorna el estado del pedido
func (o Order) Status() OrderStatus {
	return o.status
}

// CreatedAt retorna la fecha de creación
func (o Order) CreatedAt() time.Time {
	return o.createdAt
}

// UpdatedAt retorna la fecha de última actualización
func (o Order) UpdatedAt() time.Time {
	return o.updatedAt
}

// AddItem agrega un item al pedido
func (o *Order) AddItem(product Product, quantity valueobjects.Quantity) error {
	if o.status != OrderStatusPending {
		return errors.New("cannot add items to non-pending order")
	}

	// Verificar si el producto ya existe en el pedido
	for i, existingItem := range o.items {
		if existingItem.Product().Equals(product) {
			// Actualizar la cantidad del item existente
			newQuantity, err := existingItem.Quantity().Add(quantity)
			if err != nil {
				return err
			}
			return o.items[i].UpdateQuantity(newQuantity)
		}
	}

	// Crear nuevo item
	newItem, err := NewOrderItem(product, quantity)
	if err != nil {
		return err
	}

	o.items = append(o.items, *newItem)
	o.updatedAt = time.Now()
	return nil
}

// RemoveItem remueve un item del pedido
func (o *Order) RemoveItem(itemID string) error {
	if o.status != OrderStatusPending {
		return errors.New("cannot remove items from non-pending order")
	}

	for i, item := range o.items {
		if item.ID() == itemID {
			o.items = append(o.items[:i], o.items[i+1:]...)
			o.updatedAt = time.Now()
			return nil
		}
	}

	return errors.New("item not found in order")
}

// TotalAmount calcula el monto total del pedido
func (o Order) TotalAmount() (valueobjects.Money, error) {
	if len(o.items) == 0 {
		return valueobjects.Zero("USD"), nil
	}

	total := valueobjects.Zero(o.items[0].UnitPrice().Currency())

	for _, item := range o.items {
		itemTotal, err := item.TotalPrice()
		if err != nil {
			return valueobjects.Money{}, err
		}

		total, err = total.Add(itemTotal)
		if err != nil {
			return valueobjects.Money{}, err
		}
	}

	return total, nil
}

// Confirm confirma el pedido
func (o *Order) Confirm() error {
	if o.status != OrderStatusPending {
		return errors.New("only pending orders can be confirmed")
	}

	if len(o.items) == 0 {
		return errors.New("cannot confirm empty order")
	}

	o.status = OrderStatusConfirmed
	o.updatedAt = time.Now()
	return nil
}

// Ship marca el pedido como enviado
func (o *Order) Ship() error {
	if o.status != OrderStatusConfirmed {
		return errors.New("only confirmed orders can be shipped")
	}

	o.status = OrderStatusShipped
	o.updatedAt = time.Now()
	return nil
}

// Deliver marca el pedido como entregado
func (o *Order) Deliver() error {
	if o.status != OrderStatusShipped {
		return errors.New("only shipped orders can be delivered")
	}

	o.status = OrderStatusDelivered
	o.updatedAt = time.Now()
	return nil
}

// Cancel cancela el pedido
func (o *Order) Cancel() error {
	if o.status == OrderStatusDelivered {
		return errors.New("cannot cancel delivered order")
	}

	o.status = OrderStatusCancelled
	o.updatedAt = time.Now()
	return nil
}

// ItemCount retorna el número de items en el pedido
func (o Order) ItemCount() int {
	return len(o.items)
}

// IsEmpty indica si el pedido está vacío
func (o Order) IsEmpty() bool {
	return len(o.items) == 0
}

// Equals compara dos pedidos por su ID
func (o Order) Equals(other Order) bool {
	return o.id == other.id
}

// String representa el pedido como string
func (o Order) String() string {
	return fmt.Sprintf("Order %s - Customer: %s - Items: %d - Status: %s",
		o.id, o.customerID, len(o.items), o.status)
}
