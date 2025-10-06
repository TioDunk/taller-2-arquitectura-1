package services

import (
	"ddd-poc/application/dto"
	"ddd-poc/domain/entities"
	"ddd-poc/domain/repositories"
	"ddd-poc/domain/services"
	"fmt"
	"time"
)

// OrderService es un Application Service que maneja los casos de uso relacionados con pedidos
type OrderService struct {
	orderRepo         repositories.OrderRepository
	customerRepo      repositories.CustomerRepository
	productRepo       repositories.ProductRepository
	validationService *services.OrderValidationService
	pricingService    *services.PricingService
}

// NewOrderService crea una nueva instancia del servicio de pedidos
func NewOrderService(
	orderRepo repositories.OrderRepository,
	customerRepo repositories.CustomerRepository,
	productRepo repositories.ProductRepository,
) *OrderService {
	return &OrderService{
		orderRepo:         orderRepo,
		customerRepo:      customerRepo,
		productRepo:       productRepo,
		validationService: services.NewOrderValidationService(),
		pricingService:    services.NewPricingService(),
	}
}

// CreateOrder crea un nuevo pedido
func (os *OrderService) CreateOrder(request dto.CreateOrderRequest) (*dto.OrderResponse, error) {
	// Verificar que el cliente existe
	customer, err := os.customerRepo.FindByID(request.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("failed to find customer: %v", err)
	}

	if customer == nil {
		return nil, fmt.Errorf("customer with ID %s not found", request.CustomerID)
	}

	// Verificar límite de pedidos pendientes
	pendingCount, err := os.orderRepo.CountPendingByCustomer(request.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("failed to count pending orders: %v", err)
	}

	if err := os.validationService.ValidateCustomerOrderLimit(request.CustomerID, pendingCount); err != nil {
		return nil, err
	}

	// Crear la entidad Order
	order, err := entities.NewOrder(request.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %v", err)
	}

	// Guardar el pedido
	if err := os.orderRepo.Save(order); err != nil {
		return nil, fmt.Errorf("failed to save order: %v", err)
	}

	// Convertir a DTO de respuesta
	return os.orderToResponse(order), nil
}

// GetOrder obtiene un pedido por su ID
func (os *OrderService) GetOrder(id string) (*dto.OrderResponse, error) {
	order, err := os.orderRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find order: %v", err)
	}

	if order == nil {
		return nil, fmt.Errorf("order with ID %s not found", id)
	}

	return os.orderToResponse(order), nil
}

// AddItemToOrder agrega un item a un pedido
func (os *OrderService) AddItemToOrder(orderID string, request dto.AddItemRequest) (*dto.OrderResponse, error) {
	// Obtener el pedido
	order, err := os.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to find order: %v", err)
	}

	if order == nil {
		return nil, fmt.Errorf("order with ID %s not found", orderID)
	}

	// Obtener el producto
	product, err := os.productRepo.FindByID(request.ProductID)
	if err != nil {
		return nil, fmt.Errorf("failed to find product: %v", err)
	}

	if product == nil {
		return nil, fmt.Errorf("product with ID %s not found", request.ProductID)
	}

	// Convertir la cantidad
	quantity, err := request.Quantity.ToQuantity()
	if err != nil {
		return nil, fmt.Errorf("invalid quantity: %v", err)
	}

	// Validar la adición del item
	if err := os.validationService.ValidateItemAddition(*order, *product, quantity); err != nil {
		return nil, err
	}

	// Agregar el item al pedido
	if err := order.AddItem(*product, quantity); err != nil {
		return nil, fmt.Errorf("failed to add item to order: %v", err)
	}

	// Guardar los cambios
	if err := os.orderRepo.Save(order); err != nil {
		return nil, fmt.Errorf("failed to save order: %v", err)
	}

	return os.orderToResponse(order), nil
}

// ConfirmOrder confirma un pedido
func (os *OrderService) ConfirmOrder(orderID string) (*dto.OrderResponse, error) {
	// Obtener el pedido
	order, err := os.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to find order: %v", err)
	}

	if order == nil {
		return nil, fmt.Errorf("order with ID %s not found", orderID)
	}

	// Validar el pedido
	if err := os.validationService.ValidateOrder(*order); err != nil {
		return nil, fmt.Errorf("order validation failed: %v", err)
	}

	// Confirmar el pedido
	if err := order.Confirm(); err != nil {
		return nil, fmt.Errorf("failed to confirm order: %v", err)
	}

	// Guardar los cambios
	if err := os.orderRepo.Save(order); err != nil {
		return nil, fmt.Errorf("failed to save order: %v", err)
	}

	return os.orderToResponse(order), nil
}

// UpdateOrderStatus actualiza el estado de un pedido
func (os *OrderService) UpdateOrderStatus(orderID string, request dto.UpdateOrderStatusRequest) (*dto.OrderResponse, error) {
	// Obtener el pedido
	order, err := os.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to find order: %v", err)
	}

	if order == nil {
		return nil, fmt.Errorf("order with ID %s not found", orderID)
	}

	// Convertir el estado
	newStatus := entities.OrderStatus(request.Status)

	// Validar la transición
	if err := os.validationService.ValidateOrderTransition(*order, newStatus); err != nil {
		return nil, err
	}

	// Actualizar el estado según corresponda
	switch newStatus {
	case entities.OrderStatusConfirmed:
		if err := order.Confirm(); err != nil {
			return nil, fmt.Errorf("failed to confirm order: %v", err)
		}
	case entities.OrderStatusShipped:
		if err := order.Ship(); err != nil {
			return nil, fmt.Errorf("failed to ship order: %v", err)
		}
	case entities.OrderStatusDelivered:
		if err := order.Deliver(); err != nil {
			return nil, fmt.Errorf("failed to deliver order: %v", err)
		}
	case entities.OrderStatusCancelled:
		if err := order.Cancel(); err != nil {
			return nil, fmt.Errorf("failed to cancel order: %v", err)
		}
	default:
		return nil, fmt.Errorf("invalid order status: %s", request.Status)
	}

	// Guardar los cambios
	if err := os.orderRepo.Save(order); err != nil {
		return nil, fmt.Errorf("failed to save order: %v", err)
	}

	return os.orderToResponse(order), nil
}

// GetAllOrders obtiene todos los pedidos
func (os *OrderService) GetAllOrders() ([]*dto.OrderResponse, error) {
	orders, err := os.orderRepo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %v", err)
	}

	responses := make([]*dto.OrderResponse, len(orders))
	for i, order := range orders {
		responses[i] = os.orderToResponse(order)
	}

	return responses, nil
}

// orderToResponse convierte una entidad Order a OrderResponse
func (os *OrderService) orderToResponse(order *entities.Order) *dto.OrderResponse {
	// Calcular el total del pedido
	total, err := order.TotalAmount()
	if err != nil {
		// En caso de error, usar un total de 0
		total, _ = dto.MoneyRequest{Amount: 0, Currency: "USD"}.ToMoney()
	}

	// Convertir los items
	items := make([]dto.OrderItemResponse, len(order.Items()))
	for i, item := range order.Items() {
		itemTotal, _ := item.TotalPrice()
		items[i] = dto.OrderItemResponse{
			ID: item.ID(),
			Product: dto.ProductResponse{
				ID:          item.Product().ID(),
				Name:        item.Product().Name(),
				Description: item.Product().Description(),
				Price:       dto.FromMoney(item.Product().Price()),
				Stock:       dto.FromQuantity(item.Product().Stock()),
				Active:      item.Product().IsActive(),
			},
			Quantity:  dto.FromQuantity(item.Quantity()),
			UnitPrice: dto.FromMoney(item.UnitPrice()),
			Total:     dto.FromMoney(itemTotal),
		}
	}

	return &dto.OrderResponse{
		ID:         order.ID(),
		CustomerID: order.CustomerID(),
		Items:      items,
		Status:     string(order.Status()),
		Total:      dto.FromMoney(total),
		CreatedAt:  order.CreatedAt().Format(time.RFC3339),
		UpdatedAt:  order.UpdatedAt().Format(time.RFC3339),
	}
}
