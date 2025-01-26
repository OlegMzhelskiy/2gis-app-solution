package memorystore

import (
	"context"
	"testing"
	"time"

	"applicationDesignTest/internal/domain"

	"github.com/stretchr/testify/assert"
)

func TestOrderStore_AddOrder(t *testing.T) {
	tests := []struct {
		name          string
		orders        []domain.Order
		expectedOrder domain.Order
		expectedError error
	}{
		{
			name: "should add order successfully with a unique number",
			orders: []domain.Order{
				{
					ID: "1",
					Bookings: []domain.Booking{
						{HotelID: 101, RoomType: "single", From: time.Now(), To: time.Now().Add(2 * time.Hour), RoomCount: 1},
					},
				},
			},
			expectedOrder: domain.Order{
				ID:     "1",
				Number: 1,
				Bookings: []domain.Booking{
					{HotelID: 101, RoomType: "single", From: time.Now(), To: time.Now().Add(2 * time.Hour), RoomCount: 1},
				},
			},
			expectedError: nil,
		},
		{
			name: "should increment order number for each new order",
			orders: []domain.Order{
				{
					ID: "1",
					Bookings: []domain.Booking{
						{HotelID: 101, RoomType: "single", From: time.Now(), To: time.Now().Add(2 * time.Hour), RoomCount: 1},
					},
				},
				{
					ID: "2",
					Bookings: []domain.Booking{
						{HotelID: 102, RoomType: "double", From: time.Now(), To: time.Now().Add(2 * time.Hour), RoomCount: 2},
					},
				},
			},
			expectedOrder: domain.Order{
				ID:     "2",
				Number: 2,
				Bookings: []domain.Booking{
					{HotelID: 102, RoomType: "double", From: time.Now(), To: time.Now().Add(2 * time.Hour), RoomCount: 2},
				},
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewOrderStore()

			var (
				result *domain.Order
				err    error
			)

			for _, order := range tt.orders {
				result, err = store.AddOrder(context.Background(), order)
			}

			if tt.expectedError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedOrder.Number, result.Number)

			store.idMu.Lock()
			tt.expectedOrder.CreatedAt = store.ordersByID[tt.expectedOrder.ID].CreatedAt
			assert.Equal(t, tt.expectedOrder, *store.ordersByID[tt.expectedOrder.ID])
			store.idMu.Unlock()

			store.numMu.Lock()
			tt.expectedOrder.CreatedAt = store.ordersByNumber[tt.expectedOrder.Number].CreatedAt
			assert.Equal(t, tt.expectedOrder, *store.ordersByNumber[tt.expectedOrder.Number])
			store.numMu.Unlock()

			assert.Equal(t, int(store.maxOrderNumber.Load()), len(tt.orders))
		})
	}
}
