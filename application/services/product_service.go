package services

import (
	"ddd-poc/application/dto"
	"ddd-poc/domain/entities"
	"ddd-poc/domain/repositories"
	"fmt"
)

// ProductService es un Application Service que maneja los casos de uso relacionados con productos
type ProductService struct {
	productRepo repositories.ProductRepository
}

// NewProductService crea una nueva instancia del servicio de productos
func NewProductService(productRepo repositories.ProductRepository) *ProductService {
	return &ProductService{
		productRepo: productRepo,
	}
}

// CreateProduct crea un nuevo producto
func (ps *ProductService) CreateProduct(request dto.CreateProductRequest) (*dto.ProductResponse, error) {
	// Convertir el precio del DTO al value object
	price, err := request.Price.ToMoney()
	if err != nil {
		return nil, fmt.Errorf("invalid price: %v", err)
	}

	// Convertir el stock del DTO al value object
	stock, err := request.Stock.ToQuantity()
	if err != nil {
		return nil, fmt.Errorf("invalid stock: %v", err)
	}

	// Crear la entidad Product
	product, err := entities.NewProduct(request.Name, request.Description, price, stock)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %v", err)
	}

	// Guardar el producto
	if err := ps.productRepo.Save(product); err != nil {
		return nil, fmt.Errorf("failed to save product: %v", err)
	}

	// Convertir a DTO de respuesta
	return ps.productToResponse(product), nil
}

// GetProduct obtiene un producto por su ID
func (ps *ProductService) GetProduct(id string) (*dto.ProductResponse, error) {
	product, err := ps.productRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find product: %v", err)
	}

	if product == nil {
		return nil, fmt.Errorf("product with ID %s not found", id)
	}

	return ps.productToResponse(product), nil
}

// GetAllProducts obtiene todos los productos
func (ps *ProductService) GetAllProducts() ([]*dto.ProductResponse, error) {
	products, err := ps.productRepo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %v", err)
	}

	responses := make([]*dto.ProductResponse, len(products))
	for i, product := range products {
		responses[i] = ps.productToResponse(product)
	}

	return responses, nil
}

// GetActiveProducts obtiene solo los productos activos
func (ps *ProductService) GetActiveProducts() ([]*dto.ProductResponse, error) {
	products, err := ps.productRepo.FindActive()
	if err != nil {
		return nil, fmt.Errorf("failed to get active products: %v", err)
	}

	responses := make([]*dto.ProductResponse, len(products))
	for i, product := range products {
		responses[i] = ps.productToResponse(product)
	}

	return responses, nil
}

// UpdateProductPrice actualiza el precio de un producto
func (ps *ProductService) UpdateProductPrice(productID string, priceRequest dto.MoneyRequest) (*dto.ProductResponse, error) {
	// Obtener el producto
	product, err := ps.productRepo.FindByID(productID)
	if err != nil {
		return nil, fmt.Errorf("failed to find product: %v", err)
	}

	if product == nil {
		return nil, fmt.Errorf("product with ID %s not found", productID)
	}

	// Convertir el nuevo precio
	newPrice, err := priceRequest.ToMoney()
	if err != nil {
		return nil, fmt.Errorf("invalid price: %v", err)
	}

	// Actualizar el precio
	product.UpdatePrice(newPrice)

	// Guardar los cambios
	if err := ps.productRepo.Save(product); err != nil {
		return nil, fmt.Errorf("failed to update product price: %v", err)
	}

	return ps.productToResponse(product), nil
}

// UpdateProductStock actualiza el stock de un producto
func (ps *ProductService) UpdateProductStock(productID string, stockRequest dto.QuantityRequest) (*dto.ProductResponse, error) {
	// Obtener el producto
	product, err := ps.productRepo.FindByID(productID)
	if err != nil {
		return nil, fmt.Errorf("failed to find product: %v", err)
	}

	if product == nil {
		return nil, fmt.Errorf("product with ID %s not found", productID)
	}

	// Convertir el nuevo stock
	newStock, err := stockRequest.ToQuantity()
	if err != nil {
		return nil, fmt.Errorf("invalid stock: %v", err)
	}

	// Actualizar el stock
	product.UpdateStock(newStock)

	// Guardar los cambios
	if err := ps.productRepo.Save(product); err != nil {
		return nil, fmt.Errorf("failed to update product stock: %v", err)
	}

	return ps.productToResponse(product), nil
}

// DeactivateProduct desactiva un producto
func (ps *ProductService) DeactivateProduct(productID string) (*dto.ProductResponse, error) {
	// Obtener el producto
	product, err := ps.productRepo.FindByID(productID)
	if err != nil {
		return nil, fmt.Errorf("failed to find product: %v", err)
	}

	if product == nil {
		return nil, fmt.Errorf("product with ID %s not found", productID)
	}

	// Desactivar el producto
	product.Deactivate()

	// Guardar los cambios
	if err := ps.productRepo.Save(product); err != nil {
		return nil, fmt.Errorf("failed to deactivate product: %v", err)
	}

	return ps.productToResponse(product), nil
}

// productToResponse convierte una entidad Product a ProductResponse
func (ps *ProductService) productToResponse(product *entities.Product) *dto.ProductResponse {
	return &dto.ProductResponse{
		ID:          product.ID(),
		Name:        product.Name(),
		Description: product.Description(),
		Price:       dto.FromMoney(product.Price()),
		Stock:       dto.FromQuantity(product.Stock()),
		Active:      product.IsActive(),
	}
}
