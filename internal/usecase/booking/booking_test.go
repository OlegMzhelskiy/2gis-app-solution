package booking

import (
	"errors"
	"testing"
	"time"

	"applicationDesignTest/internal/domain"
	"applicationDesignTest/internal/usecase/booking/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestBookingService_CreateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHotelRepo := mocks.NewMockhotelRepository(ctrl)
	mockOrderService := mocks.NewMockorderService(ctrl)

	bs := NewBookingService(mockHotelRepo, mockOrderService)

	testOrder := domain.Order{
		ID: domain.OrderID("1-test-0"),
		Bookings: []domain.Booking{
			{HotelID: 101, RoomType: "single", From: time.Now(), To: time.Now().Add(2 * time.Hour), RoomCount: 1},
		},
	}

	tests := []struct {
		name           string
		order          domain.Order
		mockSetup      func()
		expectedResult *domain.Order
		expectedError  error
	}{
		{
			name:  "successfully create",
			order: testOrder,
			mockSetup: func() {
				mockOrderService.EXPECT().GetOrderByID(testOrder.ID).Return(nil, domain.ErrOrderNotFound)
				mockHotelRepo.EXPECT().Reserve(testOrder.Bookings).Return(nil)
				mockOrderService.EXPECT().AddOrder(testOrder).Return(&testOrder, nil)
			},
			expectedResult: &testOrder,
			expectedError:  nil,
		},
		{
			name:  "order already exists",
			order: testOrder,
			mockSetup: func() {
				mockOrderService.EXPECT().GetOrderByID(testOrder.ID).Return(&testOrder, nil)
			},
			expectedResult: nil,
			expectedError:  domain.ErrOrderAlreadyExists,
		},
		{
			name:  "getting order error",
			order: testOrder,
			mockSetup: func() {
				mockOrderService.EXPECT().GetOrderByID(testOrder.ID).Return(nil, errors.New("getting order failed"))
			},
			expectedResult: nil,
			expectedError:  errors.New("get order error"),
		},
		{
			name:  "reserve error",
			order: testOrder,
			mockSetup: func() {
				mockOrderService.EXPECT().GetOrderByID(testOrder.ID).Return(nil, domain.ErrOrderNotFound)
				mockHotelRepo.EXPECT().Reserve(testOrder.Bookings).Return(errors.New("reservation failed"))
			},
			expectedResult: nil,
			expectedError:  errors.New("reservation failed"),
		},
		{
			name:  "addition order error",
			order: testOrder,
			mockSetup: func() {
				mockOrderService.EXPECT().GetOrderByID(testOrder.ID).Return(nil, domain.ErrOrderNotFound)
				mockHotelRepo.EXPECT().Reserve(testOrder.Bookings).Return(nil)
				mockOrderService.EXPECT().AddOrder(testOrder).Return(nil, errors.New("addition order failed"))
			},
			expectedResult: nil,
			expectedError:  errors.New("addition order failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			result, err := bs.CreateOrder(tt.order)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
		})
	}
}
