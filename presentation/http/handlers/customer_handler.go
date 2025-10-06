package handlers

import (
	"ddd-poc/application/dto"
	"ddd-poc/application/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CustomerHandler maneja las peticiones HTTP relacionadas con clientes
type CustomerHandler struct {
	customerService *services.CustomerService
}

// NewCustomerHandler crea una nueva instancia del manejador de clientes
func NewCustomerHandler(customerService *services.CustomerService) *CustomerHandler {
	return &CustomerHandler{
		customerService: customerService,
	}
}

// CreateCustomer maneja la creación de un nuevo cliente
func (h *CustomerHandler) CreateCustomer(c *gin.Context) {
	var request dto.CreateCustomerRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customer, err := h.customerService.CreateCustomer(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, customer)
}

// GetCustomer maneja la obtención de un cliente por ID
func (h *CustomerHandler) GetCustomer(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "customer ID is required"})
		return
	}

	customer, err := h.customerService.GetCustomer(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, customer)
}

// GetAllCustomers maneja la obtención de todos los clientes
func (h *CustomerHandler) GetAllCustomers(c *gin.Context) {
	customers, err := h.customerService.GetAllCustomers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"customers": customers})
}

// UpdateCustomerAddress maneja la actualización de la dirección de un cliente
func (h *CustomerHandler) UpdateCustomerAddress(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "customer ID is required"})
		return
	}

	var addressRequest dto.AddressRequest
	if err := c.ShouldBindJSON(&addressRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customer, err := h.customerService.UpdateCustomerAddress(id, addressRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, customer)
}

// DeleteCustomer maneja la eliminación de un cliente
func (h *CustomerHandler) DeleteCustomer(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "customer ID is required"})
		return
	}

	// En una implementación real, esto sería manejado por el servicio
	c.JSON(http.StatusOK, gin.H{"message": "Customer deleted successfully"})
}
