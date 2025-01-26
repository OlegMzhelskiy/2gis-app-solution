package order

import "applicationDesignTest/internal/domain"

type orderRepository interface {
	GetOrderByID(id domain.OrderID) (*domain.Order, error)
	GetOrderByNumber(orderNumber domain.OrderNumber) (*domain.Order, error)
	AddOrder(order domain.Order) (*domain.Order, error)
}

type OrderService struct {
	orderStore orderRepository
}

func NewOrderService(orderStore orderRepository) *OrderService {
	return &OrderService{
		orderStore: orderStore,
	}
}

func (s *OrderService) GetOrderByID(id domain.OrderID) (*domain.Order, error) {
	return s.orderStore.GetOrderByID(id)
}

func (s *OrderService) GetOrderByNumber(orderNumber domain.OrderNumber) (*domain.Order, error) {
	return s.orderStore.GetOrderByNumber(orderNumber)
}

func (s *OrderService) AddOrder(order domain.Order) (*domain.Order, error) {
	return s.orderStore.AddOrder(order)
}
