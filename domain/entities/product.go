package entities

import (
	"ddd-poc/domain/valueobjects"
	"errors"

	"github.com/google/uuid"
)

// Product representa un producto en el catálogo
// Es una Entity porque tiene identidad única (ID)
type Product struct {
	id          string
	name        string
	description string
	price       valueobjects.Money
	stock       valueobjects.Quantity
	active      bool
}

// NewProduct crea un nuevo producto
func NewProduct(name, description string, price valueobjects.Money, stock valueobjects.Quantity) (*Product, error) {
	if name == "" {
		return nil, errors.New("product name cannot be empty")
	}
	if description == "" {
		return nil, errors.New("product description cannot be empty")
	}

	return &Product{
		id:          uuid.New().String(),
		name:        name,
		description: description,
		price:       price,
		stock:       stock,
		active:      true,
	}, nil
}

// ProductFromExisting crea un producto a partir de datos existentes
func ProductFromExisting(id, name, description string, price valueobjects.Money, stock valueobjects.Quantity, active bool) *Product {
	return &Product{
		id:          id,
		name:        name,
		description: description,
		price:       price,
		stock:       stock,
		active:      active,
	}
}

// ID retorna el ID único del producto
func (p Product) ID() string {
	return p.id
}

// Name retorna el nombre del producto
func (p Product) Name() string {
	return p.name
}

// Description retorna la descripción del producto
func (p Product) Description() string {
	return p.description
}

// Price retorna el precio del producto
func (p Product) Price() valueobjects.Money {
	return p.price
}

// Stock retorna el stock disponible
func (p Product) Stock() valueobjects.Quantity {
	return p.stock
}

// IsActive indica si el producto está activo
func (p Product) IsActive() bool {
	return p.active
}

// UpdatePrice actualiza el precio del producto
func (p *Product) UpdatePrice(newPrice valueobjects.Money) {
	p.price = newPrice
}

// UpdateStock actualiza el stock del producto
func (p *Product) UpdateStock(newStock valueobjects.Quantity) {
	p.stock = newStock
}

// ReduceStock reduce el stock del producto
func (p *Product) ReduceStock(quantity valueobjects.Quantity) error {
	if !p.IsActive() {
		return errors.New("cannot reduce stock of inactive product")
	}

	// Reducir el stock
	reducedStock, err := valueobjects.NewQuantity(p.stock.Value()-quantity.Value(), p.stock.Unit())
	if err != nil {
		return err
	}

	if reducedStock.Value() < 0 {
		return errors.New("insufficient stock")
	}

	p.stock = reducedStock
	return nil
}

// Deactivate desactiva el producto
func (p *Product) Deactivate() {
	p.active = false
}

// Equals compara dos productos por su ID
func (p Product) Equals(other Product) bool {
	return p.id == other.id
}
