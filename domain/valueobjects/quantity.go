package valueobjects

import "fmt"

// Quantity representa una cantidad
// Es un Value Object porque se identifica por su valor, no por identidad
type Quantity struct {
	value int
	unit  string
}

// NewQuantity crea una nueva cantidad
func NewQuantity(value int, unit string) (Quantity, error) {
	if value < 0 {
		return Quantity{}, fmt.Errorf("quantity cannot be negative")
	}
	if unit == "" {
		return Quantity{}, fmt.Errorf("unit cannot be empty")
	}

	return Quantity{
		value: value,
		unit:  unit,
	}, nil
}

// Value retorna el valor numérico
func (q Quantity) Value() int {
	return q.value
}

// Unit retorna la unidad
func (q Quantity) Unit() string {
	return q.unit
}

// Add suma dos cantidades (deben tener la misma unidad)
func (q Quantity) Add(other Quantity) (Quantity, error) {
	if q.unit != other.unit {
		return Quantity{}, fmt.Errorf("cannot add quantities with different units")
	}

	return NewQuantity(q.value+other.value, q.unit)
}

// Equals compara dos cantidades
func (q Quantity) Equals(other Quantity) bool {
	return q.value == other.value && q.unit == other.unit
}

// String representa la cantidad como string
func (q Quantity) String() string {
	return fmt.Sprintf("%d %s", q.value, q.unit)
}

// Zero retorna una cantidad cero
func ZeroQuantity(unit string) Quantity {
	return Quantity{
		value: 0,
		unit:  unit,
	}
}
