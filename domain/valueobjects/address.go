package valueobjects

import "fmt"

// Address representa una dirección física
// Es un Value Object porque se identifica por su contenido, no por identidad
type Address struct {
	street     string
	city       string
	state      string
	postalCode string
	country    string
}

// NewAddress crea una nueva dirección
func NewAddress(street, city, state, postalCode, country string) (Address, error) {
	if street == "" {
		return Address{}, fmt.Errorf("street cannot be empty")
	}
	if city == "" {
		return Address{}, fmt.Errorf("city cannot be empty")
	}
	if state == "" {
		return Address{}, fmt.Errorf("state cannot be empty")
	}
	if postalCode == "" {
		return Address{}, fmt.Errorf("postal code cannot be empty")
	}
	if country == "" {
		return Address{}, fmt.Errorf("country cannot be empty")
	}

	return Address{
		street:     street,
		city:       city,
		state:      state,
		postalCode: postalCode,
		country:    country,
	}, nil
}

// Street retorna la calle
func (a Address) Street() string {
	return a.street
}

// City retorna la ciudad
func (a Address) City() string {
	return a.city
}

// State retorna el estado
func (a Address) State() string {
	return a.state
}

// PostalCode retorna el código postal
func (a Address) PostalCode() string {
	return a.postalCode
}

// Country retorna el país
func (a Address) Country() string {
	return a.country
}

// Equals compara dos direcciones
func (a Address) Equals(other Address) bool {
	return a.street == other.street &&
		a.city == other.city &&
		a.state == other.state &&
		a.postalCode == other.postalCode &&
		a.country == other.country
}

// String representa la dirección como string
func (a Address) String() string {
	return fmt.Sprintf("%s, %s, %s %s, %s",
		a.street, a.city, a.state, a.postalCode, a.country)
}
