package entities

import (
	"ddd-poc/domain/valueobjects"
	"errors"
	"time"

	"github.com/google/uuid"
)

// Customer representa un cliente en el sistema
// Es una Entity porque tiene identidad única (ID)
type Customer struct {
	id        string
	name      string
	email     string
	address   valueobjects.Address
	createdAt time.Time
}

// NewCustomer crea un nuevo cliente
func NewCustomer(name, email string, address valueobjects.Address) (*Customer, error) {
	if name == "" {
		return nil, errors.New("customer name cannot be empty")
	}
	if email == "" {
		return nil, errors.New("customer email cannot be empty")
	}

	return &Customer{
		id:        uuid.New().String(),
		name:      name,
		email:     email,
		address:   address,
		createdAt: time.Now(),
	}, nil
}

// CustomerFromExisting crea un cliente a partir de datos existentes (para reconstruir desde BD)
func CustomerFromExisting(id, name, email string, address valueobjects.Address, createdAt time.Time) *Customer {
	return &Customer{
		id:        id,
		name:      name,
		email:     email,
		address:   address,
		createdAt: createdAt,
	}
}

// ID retorna el ID único del cliente
func (c Customer) ID() string {
	return c.id
}

// Name retorna el nombre del cliente
func (c Customer) Name() string {
	return c.name
}

// Email retorna el email del cliente
func (c Customer) Email() string {
	return c.email
}

// Address retorna la dirección del cliente
func (c Customer) Address() valueobjects.Address {
	return c.address
}

// CreatedAt retorna la fecha de creación
func (c Customer) CreatedAt() time.Time {
	return c.createdAt
}

// UpdateAddress actualiza la dirección del cliente
func (c *Customer) UpdateAddress(newAddress valueobjects.Address) {
	c.address = newAddress
}

// Equals compara dos clientes por su ID
func (c Customer) Equals(other Customer) bool {
	return c.id == other.id
}
