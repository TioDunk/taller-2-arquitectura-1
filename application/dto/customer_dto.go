package dto

import "ddd-poc/domain/valueobjects"

// CreateCustomerRequest representa la petición para crear un cliente
type CreateCustomerRequest struct {
	Name    string         `json:"name"`
	Email   string         `json:"email"`
	Address AddressRequest `json:"address"`
}

// AddressRequest representa la dirección en las peticiones
type AddressRequest struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
}

// CustomerResponse representa la respuesta de un cliente
type CustomerResponse struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Email     string          `json:"email"`
	Address   AddressResponse `json:"address"`
	CreatedAt string          `json:"created_at"`
}

// AddressResponse representa la dirección en las respuestas
type AddressResponse struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
}

// ToAddress convierte AddressRequest a valueobjects.Address
func (ar AddressRequest) ToAddress() (valueobjects.Address, error) {
	return valueobjects.NewAddress(ar.Street, ar.City, ar.State, ar.PostalCode, ar.Country)
}

// FromAddress convierte valueobjects.Address a AddressResponse
func FromAddress(address valueobjects.Address) AddressResponse {
	return AddressResponse{
		Street:     address.Street(),
		City:       address.City(),
		State:      address.State(),
		PostalCode: address.PostalCode(),
		Country:    address.Country(),
	}
}
