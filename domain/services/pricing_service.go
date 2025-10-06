package services

import (
	"ddd-poc/domain/entities"
	"ddd-poc/domain/valueobjects"
	"fmt"
)

// PricingService es un Domain Service que maneja la lógica de precios
// No pertenece a ninguna entidad específica, pero es lógica de dominio
type PricingService struct{}

// NewPricingService crea una nueva instancia del servicio de precios
func NewPricingService() *PricingService {
	return &PricingService{}
}

// CalculateOrderTotal calcula el total de un pedido
func (ps *PricingService) CalculateOrderTotal(order entities.Order) (valueobjects.Money, error) {
	return order.TotalAmount()
}

// CalculateItemTotal calcula el total de un item específico
func (ps *PricingService) CalculateItemTotal(item entities.OrderItem) (valueobjects.Money, error) {
	return item.TotalPrice()
}

// ApplyDiscount aplica un descuento a un monto
func (ps *PricingService) ApplyDiscount(amount valueobjects.Money, discountPercentage float64) (valueobjects.Money, error) {
	if discountPercentage < 0 || discountPercentage > 100 {
		return valueobjects.Money{}, fmt.Errorf("discount percentage must be between 0 and 100")
	}

	discountFactor := (100 - discountPercentage) / 100
	return amount.Multiply(discountFactor)
}

// CalculateTax calcula el impuesto sobre un monto
func (ps *PricingService) CalculateTax(amount valueobjects.Money, taxRate float64) (valueobjects.Money, error) {
	if taxRate < 0 {
		return valueobjects.Money{}, fmt.Errorf("tax rate cannot be negative")
	}

	return amount.Multiply(taxRate / 100)
}

// CalculateShippingCost calcula el costo de envío basado en el total del pedido
func (ps *PricingService) CalculateShippingCost(order entities.Order) (valueobjects.Money, error) {
	total, err := order.TotalAmount()
	if err != nil {
		return valueobjects.Money{}, err
	}

	// Lógica de negocio para envío
	if total.Amount() >= 100 {
		// Envío gratis para pedidos mayores a $100
		return valueobjects.Zero(total.Currency()), nil
	} else if total.Amount() >= 50 {
		// Envío de $5 para pedidos entre $50-$99
		return valueobjects.NewMoney(5.0, total.Currency())
	} else {
		// Envío de $10 para pedidos menores a $50
		return valueobjects.NewMoney(10.0, total.Currency())
	}
}
