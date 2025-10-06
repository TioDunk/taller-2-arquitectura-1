package dto

import "ddd-poc/domain/valueobjects"

// CreateProductRequest representa la petición para crear un producto
type CreateProductRequest struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Price       MoneyRequest    `json:"price"`
	Stock       QuantityRequest `json:"stock"`
}

// MoneyRequest representa el dinero en las peticiones
type MoneyRequest struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

// QuantityRequest representa la cantidad en las peticiones
type QuantityRequest struct {
	Value int    `json:"value"`
	Unit  string `json:"unit"`
}

// ProductResponse representa la respuesta de un producto
type ProductResponse struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Price       MoneyResponse    `json:"price"`
	Stock       QuantityResponse `json:"stock"`
	Active      bool             `json:"active"`
}

// MoneyResponse representa el dinero en las respuestas
type MoneyResponse struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

// QuantityResponse representa la cantidad en las respuestas
type QuantityResponse struct {
	Value int    `json:"value"`
	Unit  string `json:"unit"`
}

// ToMoney convierte MoneyRequest a valueobjects.Money
func (mr MoneyRequest) ToMoney() (valueobjects.Money, error) {
	return valueobjects.NewMoney(mr.Amount, mr.Currency)
}

// ToQuantity convierte QuantityRequest a valueobjects.Quantity
func (qr QuantityRequest) ToQuantity() (valueobjects.Quantity, error) {
	return valueobjects.NewQuantity(qr.Value, qr.Unit)
}

// FromMoney convierte valueobjects.Money a MoneyResponse
func FromMoney(money valueobjects.Money) MoneyResponse {
	return MoneyResponse{
		Amount:   money.Amount(),
		Currency: money.Currency(),
	}
}

// FromQuantity convierte valueobjects.Quantity a QuantityResponse
func FromQuantity(quantity valueobjects.Quantity) QuantityResponse {
	return QuantityResponse{
		Value: quantity.Value(),
		Unit:  quantity.Unit(),
	}
}
