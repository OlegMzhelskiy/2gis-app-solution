package memorystore

import (
	"fmt"
	"testing"
	"time"

	"applicationDesignTest/internal/domain"
	"applicationDesignTest/pkg/date"

	"github.com/stretchr/testify/assert"
)

func TestHotelStore_Reserve(t *testing.T) {
	testDate := date.Date(2025, 1, 1)

	type availability struct {
		rooms    map[time.Time]int
		roomType domain.RoomType
		hotelID  domain.HotelID
	}

	tests := []struct {
		name                 string
		roomType             domain.RoomType
		setupBookings        func(roomType domain.RoomType) []domain.Booking
		expectedError        error
		setupHotelStore      func(s *HotelStore)
		expectedAvailability []availability
	}{
		{
			name:     "successfully reserve 1 hotel",
			roomType: "single",
			setupBookings: func(roomType domain.RoomType) []domain.Booking {
				return []domain.Booking{
					{HotelID: 1, RoomType: roomType, From: testDate, To: testDate, RoomCount: 1},
				}
			},
			expectedError: nil,
			setupHotelStore: func(s *HotelStore) {
				hotel := &domain.Hotel{
					ID:   1,
					Name: "Hotel A",
				}
				roomCategory := &RoomCategory{
					availability: map[time.Time]int{
						testDate: 2,
					},
				}
				hotelWrapper := &HotelWrapper{
					Hotel: hotel,
					RoomCategories: map[domain.RoomType]*RoomCategory{
						"single": roomCategory,
					},
				}
				s.roomAvailability[1] = hotelWrapper
			},
			expectedAvailability: []availability{
				{
					rooms:    map[time.Time]int{testDate: 1},
					roomType: "single",
					hotelID:  1,
				},
			},
		},
		{
			name:     "successfully reserve 2 hotels",
			roomType: "single",
			setupBookings: func(roomType domain.RoomType) []domain.Booking {
				return []domain.Booking{
					{HotelID: 1, RoomType: roomType, From: testDate, To: testDate, RoomCount: 1},
					{HotelID: 2, RoomType: roomType, From: testDate, To: testDate.AddDate(0, 0, 1), RoomCount: 1},
				}
			},
			expectedError: nil,
			setupHotelStore: func(s *HotelStore) {
				hotel := &domain.Hotel{
					ID:   1,
					Name: "Hotel A",
				}
				roomCategory := &RoomCategory{
					availability: map[time.Time]int{
						testDate: 2,
					},
				}
				hotelWrapper := &HotelWrapper{
					Hotel: hotel,
					RoomCategories: map[domain.RoomType]*RoomCategory{
						"single": roomCategory,
					},
				}
				s.roomAvailability[1] = hotelWrapper

				hotel2 := &domain.Hotel{
					ID:   2,
					Name: "Hotel B",
				}
				roomCategory2 := &RoomCategory{
					availability: map[time.Time]int{
						testDate:                  1,
						testDate.AddDate(0, 0, 1): 2,
					},
				}
				hotelWrapper2 := &HotelWrapper{
					Hotel: hotel2,
					RoomCategories: map[domain.RoomType]*RoomCategory{
						"single": roomCategory2,
					},
				}
				s.roomAvailability[2] = hotelWrapper2
			},
			expectedAvailability: []availability{
				{
					rooms: map[time.Time]int{
						testDate: 1,
					},
					roomType: "single",
					hotelID:  1,
				},
				{
					rooms: map[time.Time]int{
						testDate:                  0,
						testDate.AddDate(0, 0, 1): 1,
					},
					roomType: "single",
					hotelID:  2,
				},
			},
		},
		{
			name:     "available only 2 of 3 rooms",
			roomType: "single",
			setupBookings: func(roomType domain.RoomType) []domain.Booking {
				return []domain.Booking{
					{HotelID: 1, RoomType: roomType, From: testDate, To: testDate, RoomCount: 3},
				}

			},
			expectedError: fmt.Errorf("room single not available in hotel id=1 for all requested dates"),
			setupHotelStore: func(s *HotelStore) {
				hotel := &domain.Hotel{
					ID:   1,
					Name: "Hotel A",
				}
				roomCategory := &RoomCategory{
					availability: map[time.Time]int{
						testDate: 2,
					},
				}
				hotelWrapper := &HotelWrapper{
					Hotel: hotel,
					RoomCategories: map[domain.RoomType]*RoomCategory{
						"single": roomCategory,
					},
				}
				s.roomAvailability[1] = hotelWrapper
			},
			expectedAvailability: []availability{
				{
					rooms: map[time.Time]int{
						testDate: 2,
					},
					roomType: "single",
					hotelID:  1,
				},
			},
		},
		{
			name:     "room type is not found",
			roomType: "double",
			setupBookings: func(roomType domain.RoomType) []domain.Booking {
				return []domain.Booking{
					{HotelID: 1, RoomType: roomType, From: testDate, To: testDate, RoomCount: 1},
				}
			},
			expectedError: domain.ErrRoomTypeNotFound,
			setupHotelStore: func(s *HotelStore) {
				hotel := &domain.Hotel{
					ID:   1,
					Name: "Hotel A",
				}
				roomCategory := &RoomCategory{
					availability: map[time.Time]int{
						testDate: 2,
					},
				}
				hotelWrapper := &HotelWrapper{
					Hotel: hotel,
					RoomCategories: map[domain.RoomType]*RoomCategory{
						"single": roomCategory,
					},
				}
				s.roomAvailability[1] = hotelWrapper
			},
			expectedAvailability: []availability{
				{
					rooms: map[time.Time]int{
						testDate: 2,
					},
					roomType: "single",
					hotelID:  1,
				},
			},
		},
		{
			name:     "available only 2 of 3 days",
			roomType: "single",
			setupBookings: func(roomType domain.RoomType) []domain.Booking {
				return []domain.Booking{
					{HotelID: 1, RoomType: roomType, From: testDate, To: testDate.AddDate(0, 0, 2), RoomCount: 3},
				}
			},
			expectedError: fmt.Errorf("room single not available in hotel id=1 for all requested dates"),
			setupHotelStore: func(s *HotelStore) {
				hotel := &domain.Hotel{
					ID:   1,
					Name: "Hotel A",
				}
				roomCategory := &RoomCategory{
					availability: map[time.Time]int{
						testDate:                  5,
						testDate.AddDate(0, 0, 1): 5,
						testDate.AddDate(0, 0, 2): 1,
					},
				}
				hotelWrapper := &HotelWrapper{
					Hotel: hotel,
					RoomCategories: map[domain.RoomType]*RoomCategory{
						"single": roomCategory,
					},
				}
				s.roomAvailability[1] = hotelWrapper
			},
			expectedAvailability: []availability{
				{
					rooms: map[time.Time]int{
						testDate:                  5,
						testDate.AddDate(0, 0, 1): 5,
						testDate.AddDate(0, 0, 2): 1,
					},
					roomType: "single",
					hotelID:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewHotelStore()

			tt.setupHotelStore(store)

			bookings := tt.setupBookings(tt.roomType)

			err := store.Reserve(bookings)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			for _, expectedAvailability := range tt.expectedAvailability {
				for date, avail := range expectedAvailability.rooms {
					hotelWrapper := store.roomAvailability[expectedAvailability.hotelID]
					category := hotelWrapper.RoomCategories[expectedAvailability.roomType]
					assert.Equal(t, avail, category.availability[date])
				}
			}
		})
	}
}
