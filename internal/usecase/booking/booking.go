package booking

//go:generate mockgen -source=booking.go -destination=mocks/mock.go -package=mocks

import (
	"context"
	"errors"
	"fmt"
	"time"

	"applicationDesignTest/internal/domain"
)

type hotelRepository interface {
	Reserve(ctx context.Context, bookings []domain.Booking) error
	AddRoomAvailability(ctx context.Context, hotelID domain.HotelID, roomType domain.RoomType, date time.Time, rooms int) error
}

type orderService interface {
	AddOrder(ctx context.Context, order domain.Order) (*domain.Order, error)
	GetOrderByID(ctx context.Context, id domain.OrderID) (*domain.Order, error)
}

type BookingService struct {
	hotelStore   hotelRepository
	orderService orderService
}

func NewBookingService(hotelStore hotelRepository, orderService orderService) *BookingService {
	return &BookingService{
		hotelStore:   hotelStore,
		orderService: orderService,
	}
}

func (bs *BookingService) CreateOrder(ctx context.Context, order domain.Order) (*domain.Order, error) {
	existOrder, err := bs.orderService.GetOrderByID(ctx, order.ID)
	if err != nil {
		if !errors.Is(err, domain.ErrOrderNotFound) {
			return nil, fmt.Errorf("failed to get order by id: %w", err)
		}
	}

	// idempotency
	if existOrder != nil {
		return existOrder, domain.ErrOrderAlreadyExists
	}

	if err := bs.hotelStore.Reserve(ctx, order.Bookings); err != nil {
		return nil, err
	}

	return bs.orderService.AddOrder(ctx, order)
}

func (bs *BookingService) AddRoomAvailability(ctx context.Context, hotelID domain.HotelID, roomType domain.RoomType, date time.Time, rooms int) error {
	return bs.hotelStore.AddRoomAvailability(ctx, hotelID, roomType, date, rooms)
}
