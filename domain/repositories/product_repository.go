package repositories

import (
	"ddd-poc/domain/entities"
	"ddd-poc/domain/valueobjects"
)

// ProductRepository define la interfaz para el repositorio de productos
type ProductRepository interface {
	// Save guarda un producto
	Save(product *entities.Product) error

	// FindByID busca un producto por su ID
	FindByID(id string) (*entities.Product, error)

	// FindByName busca productos por nombre (puede retornar múltiples)
	FindByName(name string) ([]*entities.Product, error)

	// FindAll retorna todos los productos activos
	FindAll() ([]*entities.Product, error)

	// FindActive retorna solo los productos activos
	FindActive() ([]*entities.Product, error)

	// Delete elimina un producto
	Delete(id string) error

	// Exists verifica si un producto existe
	Exists(id string) (bool, error)

	// UpdateStock actualiza el stock de un producto
	UpdateStock(productID string, newStock valueobjects.Quantity) error
}
