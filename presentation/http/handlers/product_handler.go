package handlers

import (
	"ddd-poc/application/dto"
	"ddd-poc/application/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ProductHandler maneja las peticiones HTTP relacionadas con productos
type ProductHandler struct {
	productService *services.ProductService
}

// NewProductHandler crea una nueva instancia del manejador de productos
func NewProductHandler(productService *services.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// CreateProduct maneja la creación de un nuevo producto
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var request dto.CreateProductRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.productService.CreateProduct(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// GetProduct maneja la obtención de un producto por ID
func (h *ProductHandler) GetProduct(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product ID is required"})
		return
	}

	product, err := h.productService.GetProduct(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// GetAllProducts maneja la obtención de todos los productos
func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	products, err := h.productService.GetAllProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"products": products})
}

// GetActiveProducts maneja la obtención de productos activos
func (h *ProductHandler) GetActiveProducts(c *gin.Context) {
	products, err := h.productService.GetActiveProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"products": products})
}

// UpdateProductPrice maneja la actualización del precio de un producto
func (h *ProductHandler) UpdateProductPrice(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product ID is required"})
		return
	}

	var priceRequest dto.MoneyRequest
	if err := c.ShouldBindJSON(&priceRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.productService.UpdateProductPrice(id, priceRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// UpdateProductStock maneja la actualización del stock de un producto
func (h *ProductHandler) UpdateProductStock(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product ID is required"})
		return
	}

	var stockRequest dto.QuantityRequest
	if err := c.ShouldBindJSON(&stockRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.productService.UpdateProductStock(id, stockRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// DeactivateProduct maneja la desactivación de un producto
func (h *ProductHandler) DeactivateProduct(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product ID is required"})
		return
	}

	product, err := h.productService.DeactivateProduct(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}
