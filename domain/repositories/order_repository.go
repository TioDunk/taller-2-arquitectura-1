package repositories

import "ddd-poc/domain/entities"

// OrderRepository define la interfaz para el repositorio de pedidos
type OrderRepository interface {
	// Save guarda un pedido
	Save(order *entities.Order) error

	// FindByID busca un pedido por su ID
	FindByID(id string) (*entities.Order, error)

	// FindByCustomerID busca pedidos por ID de cliente
	FindByCustomerID(customerID string) ([]*entities.Order, error)

	// FindByStatus busca pedidos por estado
	FindByStatus(status entities.OrderStatus) ([]*entities.Order, error)

	// FindAll retorna todos los pedidos
	FindAll() ([]*entities.Order, error)

	// Delete elimina un pedido
	Delete(id string) error

	// Exists verifica si un pedido existe
	Exists(id string) (bool, error)

	// CountPendingByCustomer cuenta los pedidos pendientes de un cliente
	CountPendingByCustomer(customerID string) (int, error)
}
