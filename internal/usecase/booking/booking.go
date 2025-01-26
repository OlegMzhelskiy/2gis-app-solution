package booking

//go:generate mockgen -source=booking.go -destination=mocks/mock.go -package=mocks

import (
	"errors"
	"fmt"
	"time"

	"applicationDesignTest/internal/domain"
)

type hotelRepository interface {
	Reserve(bookings []domain.Booking) error
	AddRoomAvailability(hotelID domain.HotelID, roomType domain.RoomType, date time.Time, rooms int) error
}

type orderService interface {
	AddOrder(order domain.Order) (*domain.Order, error)
	GetOrderByID(id domain.OrderID) (*domain.Order, error)
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

func (bs *BookingService) CreateOrder(order domain.Order) (*domain.Order, error) {
	existOrder, err := bs.orderService.GetOrderByID(order.ID)
	if err != nil {
		if !errors.Is(err, domain.ErrOrderNotFound) {
			return nil, fmt.Errorf("failed to get order by id: %w", err)
		}
	}

	// idempotency
	if existOrder != nil {
		return existOrder, domain.ErrOrderAlreadyExists
	}

	if err := bs.hotelStore.Reserve(order.Bookings); err != nil {
		return nil, err
	}

	return bs.orderService.AddOrder(order)
}

func (bs *BookingService) AddRoomAvailability(hotelID domain.HotelID, roomType domain.RoomType, date time.Time, rooms int) error {
	return bs.hotelStore.AddRoomAvailability(hotelID, roomType, date, rooms)
}
