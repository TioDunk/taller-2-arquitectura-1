package database

import (
	"ddd-poc/domain/entities"
	"sync"
)

// InMemoryDatabase es una base de datos en memoria para la POC
// En un proyecto real, esto sería reemplazado por una base de datos real
type InMemoryDatabase struct {
	customers map[string]*entities.Customer
	products  map[string]*entities.Product
	orders    map[string]*entities.Order

	// Mutex para operaciones concurrentes
	customerMutex sync.RWMutex
	productMutex  sync.RWMutex
	orderMutex    sync.RWMutex
}

// NewInMemoryDatabase crea una nueva instancia de la base de datos en memoria
func NewInMemoryDatabase() *InMemoryDatabase {
	return &InMemoryDatabase{
		customers: make(map[string]*entities.Customer),
		products:  make(map[string]*entities.Product),
		orders:    make(map[string]*entities.Order),
	}
}

// GetCustomers retorna el mapa de clientes (para uso interno de los repositorios)
func (db *InMemoryDatabase) GetCustomers() map[string]*entities.Customer {
	return db.customers
}

// GetCustomerMutex retorna el mutex de clientes
func (db *InMemoryDatabase) GetCustomerMutex() *sync.RWMutex {
	return &db.customerMutex
}

// GetProducts retorna el mapa de productos
func (db *InMemoryDatabase) GetProducts() map[string]*entities.Product {
	return db.products
}

// GetProductMutex retorna el mutex de productos
func (db *InMemoryDatabase) GetProductMutex() *sync.RWMutex {
	return &db.productMutex
}

// GetOrders retorna el mapa de pedidos
func (db *InMemoryDatabase) GetOrders() map[string]*entities.Order {
	return db.orders
}

// GetOrderMutex retorna el mutex de pedidos
func (db *InMemoryDatabase) GetOrderMutex() *sync.RWMutex {
	return &db.orderMutex
}
