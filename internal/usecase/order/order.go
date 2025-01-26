package order

import (
	"context"

	"applicationDesignTest/internal/domain"
)

type orderRepository interface {
	GetOrderByID(ctx context.Context, id domain.OrderID) (*domain.Order, error)
	GetOrderByNumber(ctx context.Context, orderNumber domain.OrderNumber) (*domain.Order, error)
	AddOrder(ctx context.Context, order domain.Order) (*domain.Order, error)
}

type OrderService struct {
	orderStore orderRepository
}

func NewOrderService(orderStore orderRepository) *OrderService {
	return &OrderService{
		orderStore: orderStore,
	}
}

func (s *OrderService) GetOrderByID(ctx context.Context, id domain.OrderID) (*domain.Order, error) {
	return s.orderStore.GetOrderByID(ctx, id)
}

func (s *OrderService) GetOrderByNumber(ctx context.Context, orderNumber domain.OrderNumber) (*domain.Order, error) {
	return s.orderStore.GetOrderByNumber(ctx, orderNumber)
}

func (s *OrderService) AddOrder(ctx context.Context, order domain.Order) (*domain.Order, error) {
	return s.orderStore.AddOrder(ctx, order)
}
