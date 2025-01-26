package memorystore

//go:generate mockgen -source=hotel.go -destination=mocks/mock.go -package=mocks

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"applicationDesignTest/internal/domain"
)

type HotelStore struct {
	roomAvailability map[domain.HotelID]*HotelWrapper
	mu               sync.RWMutex // lock for addition new hotel
}

type HotelWrapper struct {
	Hotel          *domain.Hotel
	RoomCategories map[domain.RoomType]*RoomCategory // RoomType -> RoomCategory
	mu             sync.Mutex
}

type RoomCategory struct {
	availability map[time.Time]int // Date -> Available Rooms
	mu           sync.Mutex
}

type reservedCategories struct {
	category  *RoomCategory
	from      time.Time
	to        time.Time
	roomCount int
}

func NewHotelStore() *HotelStore {
	return &HotelStore{
		roomAvailability: make(map[domain.HotelID]*HotelWrapper),
	}
}

func (s *HotelStore) GetHotel(ctx context.Context, hotelID domain.HotelID) (*domain.Hotel, error) {
	s.mu.RLock()
	hotelWrapper, ok := s.roomAvailability[hotelID]
	s.mu.RUnlock()

	if !ok {
		return nil, domain.ErrHotelNotFound
	}

	return hotelWrapper.Hotel, nil
}

func (s *HotelStore) AddHotel(ctx context.Context, hotel domain.Hotel) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.roomAvailability[hotel.ID]; !ok {
		s.roomAvailability[hotel.ID] = &HotelWrapper{
			Hotel:          &hotel,
			RoomCategories: make(map[domain.RoomType]*RoomCategory),
		}
	}

	return nil
}

func (s *HotelStore) AddRoomAvailability(ctx context.Context, hotelID domain.HotelID, roomType domain.RoomType, date time.Time, rooms int) error {
	s.mu.RLock()
	hotelWrapper, ok := s.roomAvailability[hotelID]
	s.mu.RUnlock()

	if !ok {
		return domain.ErrHotelNotFound
	}

	roomCat, ok := hotelWrapper.RoomCategories[roomType]
	if !ok {
		hotelWrapper.mu.Lock()
		hotelWrapper.RoomCategories[roomType] = &RoomCategory{
			availability: map[time.Time]int{date: rooms},
		}
		hotelWrapper.mu.Unlock()

		return nil
	}

	roomCat.mu.Lock()
	roomCat.availability[date] += rooms
	roomCat.mu.Unlock()

	return nil
}

func (s *HotelStore) Reserve(ctx context.Context, bookings []domain.Booking) error {
	var lockedCategories []reservedCategories

	defer func() {
		for _, reserve := range lockedCategories {
			reserve.category.mu.Unlock()
		}
	}()

	// resolving deadlocks
	sort.Slice(bookings, func(i, j int) bool {
		return bookings[i].HotelID < bookings[j].HotelID
	})

	// checking availability
	for _, booking := range bookings {
		s.mu.RLock()
		hotelWrapper, ok := s.roomAvailability[booking.HotelID]
		s.mu.RUnlock()

		if !ok {
			return domain.ErrHotelNotFound
		}

		category, ok := hotelWrapper.RoomCategories[booking.RoomType]
		if !ok {
			return domain.ErrRoomTypeNotFound
		}

		category.mu.Lock()

		lockedCategories = append(lockedCategories, reservedCategories{
			category:  category,
			from:      booking.From,
			to:        booking.To,
			roomCount: booking.RoomCount,
		})

		for date := booking.From; !date.After(booking.To); date = date.AddDate(0, 0, 1) {
			if category.availability[date] < booking.RoomCount {
				return fmt.Errorf("%w: room '%s' not available in hotel id=%v for all requested dates",
					domain.ErrRoomsNotAvailable, booking.RoomType, booking.HotelID)
			}
		}
	}

	// decrease availability
	for _, reserve := range lockedCategories {
		for date := reserve.from; !date.After(reserve.to); date = date.AddDate(0, 0, 1) {
			reserve.category.availability[date] -= reserve.roomCount
		}
	}

	return nil
}
