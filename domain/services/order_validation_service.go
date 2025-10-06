package services

import (
	"ddd-poc/domain/entities"
	"ddd-poc/domain/valueobjects"
	"fmt"
)

// OrderValidationService es un Domain Service que maneja validaciones de pedidos
type OrderValidationService struct {
	pricingService *PricingService
}

// NewOrderValidationService crea una nueva instancia del servicio de validación
func NewOrderValidationService() *OrderValidationService {
	return &OrderValidationService{
		pricingService: NewPricingService(),
	}
}

// ValidateOrder valida un pedido antes de confirmarlo
func (ovs *OrderValidationService) ValidateOrder(order entities.Order) error {
	// Validar que el pedido no esté vacío
	if order.IsEmpty() {
		return fmt.Errorf("order cannot be empty")
	}

	// Validar que todos los productos estén activos
	for _, item := range order.Items() {
		if !item.Product().IsActive() {
			return fmt.Errorf("cannot include inactive product: %s", item.Product().Name())
		}
	}

	// Validar stock disponible (simplificado)
	for _, item := range order.Items() {
		if item.Quantity().Value() > item.Product().Stock().Value() {
			return fmt.Errorf("insufficient stock for product: %s", item.Product().Name())
		}
	}

	// Validar monto mínimo del pedido
	total, err := ovs.pricingService.CalculateOrderTotal(order)
	if err != nil {
		return fmt.Errorf("error calculating order total: %v", err)
	}

	if total.Amount() < 1.0 {
		return fmt.Errorf("order total must be at least $1.00")
	}

	return nil
}

// ValidateItemAddition valida si se puede agregar un item al pedido
func (ovs *OrderValidationService) ValidateItemAddition(order entities.Order, product entities.Product, quantity valueobjects.Quantity) error {
	// Validar que el producto esté activo
	if !product.IsActive() {
		return fmt.Errorf("cannot add inactive product: %s", product.Name())
	}

	// Validar cantidad
	if quantity.Value() <= 0 {
		return fmt.Errorf("quantity must be greater than zero")
	}

	// Validar stock disponible
	if quantity.Value() > product.Stock().Value() {
		return fmt.Errorf("insufficient stock for product: %s (available: %d, requested: %d)",
			product.Name(), product.Stock().Value(), quantity.Value())
	}

	// Validar que el pedido esté en estado pendiente
	if order.Status() != entities.OrderStatusPending {
		return fmt.Errorf("cannot add items to order in status: %s", order.Status())
	}

	return nil
}

// ValidateOrderTransition valida si se puede hacer una transición de estado
func (ovs *OrderValidationService) ValidateOrderTransition(order entities.Order, newStatus entities.OrderStatus) error {
	currentStatus := order.Status()

	// Definir transiciones válidas
	validTransitions := map[entities.OrderStatus][]entities.OrderStatus{
		entities.OrderStatusPending:   {entities.OrderStatusConfirmed, entities.OrderStatusCancelled},
		entities.OrderStatusConfirmed: {entities.OrderStatusShipped, entities.OrderStatusCancelled},
		entities.OrderStatusShipped:   {entities.OrderStatusDelivered},
		entities.OrderStatusDelivered: {}, // Estado final
		entities.OrderStatusCancelled: {}, // Estado final
	}

	allowedStatuses, exists := validTransitions[currentStatus]
	if !exists {
		return fmt.Errorf("invalid current status: %s", currentStatus)
	}

	for _, allowedStatus := range allowedStatuses {
		if allowedStatus == newStatus {
			return nil
		}
	}

	return fmt.Errorf("invalid transition from %s to %s", currentStatus, newStatus)
}

// ValidateCustomerOrderLimit valida límites de pedidos por cliente (ejemplo)
func (ovs *OrderValidationService) ValidateCustomerOrderLimit(customerID string, existingOrdersCount int) error {
	// Límite de ejemplo: máximo 10 pedidos pendientes por cliente
	maxPendingOrders := 10

	if existingOrdersCount >= maxPendingOrders {
		return fmt.Errorf("customer %s has reached the maximum limit of pending orders (%d)",
			customerID, maxPendingOrders)
	}

	return nil
}
