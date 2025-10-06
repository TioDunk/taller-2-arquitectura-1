package repositories

import (
	"ddd-poc/domain/entities"
	"ddd-poc/infrastructure/database"
	"fmt"
)

// InMemoryCustomerRepository implementa CustomerRepository usando base de datos en memoria
type InMemoryCustomerRepository struct {
	db *database.InMemoryDatabase
}

// NewInMemoryCustomerRepository crea una nueva instancia del repositorio en memoria
func NewInMemoryCustomerRepository(db *database.InMemoryDatabase) *InMemoryCustomerRepository {
	return &InMemoryCustomerRepository{
		db: db,
	}
}

// Save guarda un cliente
func (r *InMemoryCustomerRepository) Save(customer *entities.Customer) error {
	r.db.GetCustomerMutex().Lock()
	defer r.db.GetCustomerMutex().Unlock()

	r.db.GetCustomers()[customer.ID()] = customer
	return nil
}

// FindByID busca un cliente por su ID
func (r *InMemoryCustomerRepository) FindByID(id string) (*entities.Customer, error) {
	r.db.GetCustomerMutex().RLock()
	defer r.db.GetCustomerMutex().RUnlock()

	customer, exists := r.db.GetCustomers()[id]
	if !exists {
		return nil, nil
	}

	return customer, nil
}

// FindByEmail busca un cliente por su email
func (r *InMemoryCustomerRepository) FindByEmail(email string) (*entities.Customer, error) {
	r.db.GetCustomerMutex().RLock()
	defer r.db.GetCustomerMutex().RUnlock()

	for _, customer := range r.db.GetCustomers() {
		if customer.Email() == email {
			return customer, nil
		}
	}

	return nil, nil
}

// FindAll retorna todos los clientes
func (r *InMemoryCustomerRepository) FindAll() ([]*entities.Customer, error) {
	r.db.GetCustomerMutex().RLock()
	defer r.db.GetCustomerMutex().RUnlock()

	customers := make([]*entities.Customer, 0, len(r.db.GetCustomers()))
	for _, customer := range r.db.GetCustomers() {
		customers = append(customers, customer)
	}

	return customers, nil
}

// Delete elimina un cliente
func (r *InMemoryCustomerRepository) Delete(id string) error {
	r.db.GetCustomerMutex().Lock()
	defer r.db.GetCustomerMutex().Unlock()

	if _, exists := r.db.GetCustomers()[id]; !exists {
		return fmt.Errorf("customer with ID %s not found", id)
	}

	delete(r.db.GetCustomers(), id)
	return nil
}

// Exists verifica si un cliente existe
func (r *InMemoryCustomerRepository) Exists(id string) (bool, error) {
	r.db.GetCustomerMutex().RLock()
	defer r.db.GetCustomerMutex().RUnlock()

	_, exists := r.db.GetCustomers()[id]
	return exists, nil
}
