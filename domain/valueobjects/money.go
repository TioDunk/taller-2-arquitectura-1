package valueobjects

import (
	"fmt"
	"math"
)

// Money representa un valor monetario
// Es un Value Object porque:
// 1. No tiene identidad (se compara por valor)
// 2. Es inmutable
// 3. Encapsula la lógica de manejo de dinero
type Money struct {
	amount   float64
	currency string
}

// NewMoney crea una nueva instancia de Money
func NewMoney(amount float64, currency string) (Money, error) {
	if amount < 0 {
		return Money{}, fmt.Errorf("amount cannot be negative")
	}
	if currency == "" {
		return Money{}, fmt.Errorf("currency cannot be empty")
	}

	// Redondear a 2 decimales para evitar problemas de precisión
	roundedAmount := math.Round(amount*100) / 100

	return Money{
		amount:   roundedAmount,
		currency: currency,
	}, nil
}

// Amount retorna el monto
func (m Money) Amount() float64 {
	return m.amount
}

// Currency retorna la moneda
func (m Money) Currency() string {
	return m.currency
}

// Add suma dos valores monetarios (deben tener la misma moneda)
func (m Money) Add(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, fmt.Errorf("cannot add money with different currencies")
	}

	return NewMoney(m.amount+other.amount, m.currency)
}

// Multiply multiplica el dinero por un factor
func (m Money) Multiply(factor float64) (Money, error) {
	if factor < 0 {
		return Money{}, fmt.Errorf("factor cannot be negative")
	}

	return NewMoney(m.amount*factor, m.currency)
}

// Equals compara dos valores monetarios
func (m Money) Equals(other Money) bool {
	return m.amount == other.amount && m.currency == other.currency
}

// String representa el dinero como string
func (m Money) String() string {
	return fmt.Sprintf("%.2f %s", m.amount, m.currency)
}

// Zero retorna un valor monetario cero
func Zero(currency string) Money {
	return Money{
		amount:   0.0,
		currency: currency,
	}
}
