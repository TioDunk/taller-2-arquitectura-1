package repositories

import "ddd-poc/domain/entities"

// CustomerRepository define la interfaz para el repositorio de clientes
// Esta interfaz está en la capa de dominio, no en infraestructura
// Siguiendo el principio de inversión de dependencias
type CustomerRepository interface {
	// Save guarda un cliente
	Save(customer *entities.Customer) error

	// FindByID busca un cliente por su ID
	FindByID(id string) (*entities.Customer, error)

	// FindByEmail busca un cliente por su email
	FindByEmail(email string) (*entities.Customer, error)

	// FindAll retorna todos los clientes
	FindAll() ([]*entities.Customer, error)

	// Delete elimina un cliente
	Delete(id string) error

	// Exists verifica si un cliente existe
	Exists(id string) (bool, error)
}
