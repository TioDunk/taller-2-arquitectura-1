package repositories

import (
	"ddd-poc/domain/entities"
	"ddd-poc/infrastructure/database"
	"fmt"
)

// InMemoryOrderRepository implementa OrderRepository usando base de datos en memoria
type InMemoryOrderRepository struct {
	db *database.InMemoryDatabase
}

// NewInMemoryOrderRepository crea una nueva instancia del repositorio en memoria
func NewInMemoryOrderRepository(db *database.InMemoryDatabase) *InMemoryOrderRepository {
	return &InMemoryOrderRepository{
		db: db,
	}
}

// Save guarda un pedido
func (r *InMemoryOrderRepository) Save(order *entities.Order) error {
	r.db.GetOrderMutex().Lock()
	defer r.db.GetOrderMutex().Unlock()

	r.db.GetOrders()[order.ID()] = order
	return nil
}

// FindByID busca un pedido por su ID
func (r *InMemoryOrderRepository) FindByID(id string) (*entities.Order, error) {
	r.db.GetOrderMutex().RLock()
	defer r.db.GetOrderMutex().RUnlock()

	order, exists := r.db.GetOrders()[id]
	if !exists {
		return nil, nil
	}

	return order, nil
}

// FindByCustomerID busca pedidos por ID de cliente
func (r *InMemoryOrderRepository) FindByCustomerID(customerID string) ([]*entities.Order, error) {
	r.db.GetOrderMutex().RLock()
	defer r.db.GetOrderMutex().RUnlock()

	var orders []*entities.Order
	for _, order := range r.db.GetOrders() {
		if order.CustomerID() == customerID {
			orders = append(orders, order)
		}
	}

	return orders, nil
}

// FindByStatus busca pedidos por estado
func (r *InMemoryOrderRepository) FindByStatus(status entities.OrderStatus) ([]*entities.Order, error) {
	r.db.GetOrderMutex().RLock()
	defer r.db.GetOrderMutex().RUnlock()

	var orders []*entities.Order
	for _, order := range r.db.GetOrders() {
		if order.Status() == status {
			orders = append(orders, order)
		}
	}

	return orders, nil
}

// FindAll retorna todos los pedidos
func (r *InMemoryOrderRepository) FindAll() ([]*entities.Order, error) {
	r.db.GetOrderMutex().RLock()
	defer r.db.GetOrderMutex().RUnlock()

	orders := make([]*entities.Order, 0, len(r.db.GetOrders()))
	for _, order := range r.db.GetOrders() {
		orders = append(orders, order)
	}

	return orders, nil
}

// Delete elimina un pedido
func (r *InMemoryOrderRepository) Delete(id string) error {
	r.db.GetOrderMutex().Lock()
	defer r.db.GetOrderMutex().Unlock()

	if _, exists := r.db.GetOrders()[id]; !exists {
		return fmt.Errorf("order with ID %s not found", id)
	}

	delete(r.db.GetOrders(), id)
	return nil
}

// Exists verifica si un pedido existe
func (r *InMemoryOrderRepository) Exists(id string) (bool, error) {
	r.db.GetOrderMutex().RLock()
	defer r.db.GetOrderMutex().RUnlock()

	_, exists := r.db.GetOrders()[id]
	return exists, nil
}

// CountPendingByCustomer cuenta los pedidos pendientes de un cliente
func (r *InMemoryOrderRepository) CountPendingByCustomer(customerID string) (int, error) {
	r.db.GetOrderMutex().RLock()
	defer r.db.GetOrderMutex().RUnlock()

	count := 0
	for _, order := range r.db.GetOrders() {
		if order.CustomerID() == customerID && order.Status() == entities.OrderStatusPending {
			count++
		}
	}

	return count, nil
}
