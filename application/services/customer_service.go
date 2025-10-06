package services

import (
	"ddd-poc/application/dto"
	"ddd-poc/domain/entities"
	"ddd-poc/domain/repositories"
	"fmt"
	"time"
)

// CustomerService es un Application Service que maneja los casos de uso relacionados con clientes
type CustomerService struct {
	customerRepo repositories.CustomerRepository
}

// NewCustomerService crea una nueva instancia del servicio de clientes
func NewCustomerService(customerRepo repositories.CustomerRepository) *CustomerService {
	return &CustomerService{
		customerRepo: customerRepo,
	}
}

// CreateCustomer crea un nuevo cliente
func (cs *CustomerService) CreateCustomer(request dto.CreateCustomerRequest) (*dto.CustomerResponse, error) {
	// Convertir la dirección del DTO al value object
	address, err := request.Address.ToAddress()
	if err != nil {
		return nil, fmt.Errorf("invalid address: %v", err)
	}

	// Verificar si ya existe un cliente con ese email
	existingCustomer, err := cs.customerRepo.FindByEmail(request.Email)
	if err == nil && existingCustomer != nil {
		return nil, fmt.Errorf("customer with email %s already exists", request.Email)
	}

	// Crear la entidad Customer
	customer, err := entities.NewCustomer(request.Name, request.Email, address)
	if err != nil {
		return nil, fmt.Errorf("failed to create customer: %v", err)
	}

	// Guardar el cliente
	if err := cs.customerRepo.Save(customer); err != nil {
		return nil, fmt.Errorf("failed to save customer: %v", err)
	}

	// Convertir a DTO de respuesta
	return cs.customerToResponse(customer), nil
}

// GetCustomer obtiene un cliente por su ID
func (cs *CustomerService) GetCustomer(id string) (*dto.CustomerResponse, error) {
	customer, err := cs.customerRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find customer: %v", err)
	}

	if customer == nil {
		return nil, fmt.Errorf("customer with ID %s not found", id)
	}

	return cs.customerToResponse(customer), nil
}

// GetAllCustomers obtiene todos los clientes
func (cs *CustomerService) GetAllCustomers() ([]*dto.CustomerResponse, error) {
	customers, err := cs.customerRepo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get customers: %v", err)
	}

	responses := make([]*dto.CustomerResponse, len(customers))
	for i, customer := range customers {
		responses[i] = cs.customerToResponse(customer)
	}

	return responses, nil
}

// UpdateCustomerAddress actualiza la dirección de un cliente
func (cs *CustomerService) UpdateCustomerAddress(customerID string, addressRequest dto.AddressRequest) (*dto.CustomerResponse, error) {
	// Obtener el cliente
	customer, err := cs.customerRepo.FindByID(customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to find customer: %v", err)
	}

	if customer == nil {
		return nil, fmt.Errorf("customer with ID %s not found", customerID)
	}

	// Convertir la nueva dirección
	newAddress, err := addressRequest.ToAddress()
	if err != nil {
		return nil, fmt.Errorf("invalid address: %v", err)
	}

	// Actualizar la dirección
	customer.UpdateAddress(newAddress)

	// Guardar los cambios
	if err := cs.customerRepo.Save(customer); err != nil {
		return nil, fmt.Errorf("failed to update customer address: %v", err)
	}

	return cs.customerToResponse(customer), nil
}

// customerToResponse convierte una entidad Customer a CustomerResponse
func (cs *CustomerService) customerToResponse(customer *entities.Customer) *dto.CustomerResponse {
	return &dto.CustomerResponse{
		ID:        customer.ID(),
		Name:      customer.Name(),
		Email:     customer.Email(),
		Address:   dto.FromAddress(customer.Address()),
		CreatedAt: customer.CreatedAt().Format(time.RFC3339),
	}
}
