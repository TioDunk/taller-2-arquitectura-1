package repositories

import (
	"ddd-poc/domain/entities"
	"ddd-poc/domain/valueobjects"
	"ddd-poc/infrastructure/database"
	"fmt"
	"strings"
)

// InMemoryProductRepository implementa ProductRepository usando base de datos en memoria
type InMemoryProductRepository struct {
	db *database.InMemoryDatabase
}

// NewInMemoryProductRepository crea una nueva instancia del repositorio en memoria
func NewInMemoryProductRepository(db *database.InMemoryDatabase) *InMemoryProductRepository {
	return &InMemoryProductRepository{
		db: db,
	}
}

// Save guarda un producto
func (r *InMemoryProductRepository) Save(product *entities.Product) error {
	r.db.GetProductMutex().Lock()
	defer r.db.GetProductMutex().Unlock()

	r.db.GetProducts()[product.ID()] = product
	return nil
}

// FindByID busca un producto por su ID
func (r *InMemoryProductRepository) FindByID(id string) (*entities.Product, error) {
	r.db.GetProductMutex().RLock()
	defer r.db.GetProductMutex().RUnlock()

	product, exists := r.db.GetProducts()[id]
	if !exists {
		return nil, nil
	}

	return product, nil
}

// FindByName busca productos por nombre (puede retornar múltiples)
func (r *InMemoryProductRepository) FindByName(name string) ([]*entities.Product, error) {
	r.db.GetProductMutex().RLock()
	defer r.db.GetProductMutex().RUnlock()

	var products []*entities.Product
	lowerName := strings.ToLower(name)

	for _, product := range r.db.GetProducts() {
		if strings.Contains(strings.ToLower(product.Name()), lowerName) {
			products = append(products, product)
		}
	}

	return products, nil
}

// FindAll retorna todos los productos
func (r *InMemoryProductRepository) FindAll() ([]*entities.Product, error) {
	r.db.GetProductMutex().RLock()
	defer r.db.GetProductMutex().RUnlock()

	products := make([]*entities.Product, 0, len(r.db.GetProducts()))
	for _, product := range r.db.GetProducts() {
		products = append(products, product)
	}

	return products, nil
}

// FindActive retorna solo los productos activos
func (r *InMemoryProductRepository) FindActive() ([]*entities.Product, error) {
	r.db.GetProductMutex().RLock()
	defer r.db.GetProductMutex().RUnlock()

	var products []*entities.Product
	for _, product := range r.db.GetProducts() {
		if product.IsActive() {
			products = append(products, product)
		}
	}

	return products, nil
}

// Delete elimina un producto
func (r *InMemoryProductRepository) Delete(id string) error {
	r.db.GetProductMutex().Lock()
	defer r.db.GetProductMutex().Unlock()

	if _, exists := r.db.GetProducts()[id]; !exists {
		return fmt.Errorf("product with ID %s not found", id)
	}

	delete(r.db.GetProducts(), id)
	return nil
}

// Exists verifica si un producto existe
func (r *InMemoryProductRepository) Exists(id string) (bool, error) {
	r.db.GetProductMutex().RLock()
	defer r.db.GetProductMutex().RUnlock()

	_, exists := r.db.GetProducts()[id]
	return exists, nil
}

// UpdateStock actualiza el stock de un producto
func (r *InMemoryProductRepository) UpdateStock(productID string, newStock valueobjects.Quantity) error {
	r.db.GetProductMutex().Lock()
	defer r.db.GetProductMutex().Unlock()

	product, exists := r.db.GetProducts()[productID]
	if !exists {
		return fmt.Errorf("product with ID %s not found", productID)
	}

	product.UpdateStock(newStock)
	return nil
}
