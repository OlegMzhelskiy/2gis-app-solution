package memorystore

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"applicationDesignTest/internal/domain"
)

type OrderStore struct {
	ordersByID map[domain.OrderID]*domain.Order
	idMu       sync.RWMutex

	ordersByNumber map[domain.OrderNumber]*domain.Order
	numMu          sync.RWMutex

	maxOrderNumber atomic.Int64
}

func NewOrderStore() *OrderStore {
	return &OrderStore{
		ordersByID:     make(map[domain.OrderID]*domain.Order),
		ordersByNumber: make(map[domain.OrderNumber]*domain.Order),
	}
}

func (s *OrderStore) AddOrder(ctx context.Context, order domain.Order) (*domain.Order, error) {
	order.Number = domain.OrderNumber(s.maxOrderNumber.Add(1))
	order.CreatedAt = time.Now()

	s.idMu.Lock()
	defer s.idMu.Unlock()

	s.numMu.Lock()
	defer s.numMu.Unlock()

	s.ordersByID[order.ID] = &order
	s.ordersByNumber[order.Number] = &order

	return &order, nil
}

func (s *OrderStore) GetOrderByID(ctx context.Context, id domain.OrderID) (*domain.Order, error) {
	s.idMu.RLock()
	defer s.idMu.RUnlock()

	order, ok := s.ordersByID[id]
	if !ok {
		return nil, domain.ErrOrderNotFound
	}

	return order, nil
}

func (s *OrderStore) GetOrderByNumber(ctx context.Context, orderNumber domain.OrderNumber) (*domain.Order, error) {
	s.numMu.RLock()
	defer s.numMu.RUnlock()

	order, ok := s.ordersByNumber[orderNumber]
	if !ok {
		return nil, domain.ErrOrderNotFound
	}

	return order, nil
}
