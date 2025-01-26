package create_order

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"applicationDesignTest/internal/api/http_helpers"
	"applicationDesignTest/internal/domain"
	"applicationDesignTest/pkg/date"
	"applicationDesignTest/pkg/log"
)

type request struct {
	ID       domain.OrderID `json:"id"`
	UserID   domain.UserID  `json:"user_id"`
	Bookings []booking      `json:"booking"`
}

type booking struct {
	HotelID   domain.HotelID  `json:"hotel_id"`
	RoomType  domain.RoomType `json:"room_type"`
	From      date.CustomDate `json:"from"`
	To        date.CustomDate `json:"to"`
	RoomCount int             `json:"room_count"`
}

type bookingService interface {
	CreateOrder(ctx context.Context, order domain.Order) (*domain.Order, error)
}

type Handler struct {
	booking bookingService
}

func NewHandler(bookingService bookingService) *Handler {
	return &Handler{
		booking: bookingService,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req request

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Warning(fmt.Sprintf("failed to decode request: %s", err.Error()))
		http_helpers.SendError(w, http.StatusBadRequest, "invalid input", http_helpers.ErrorTypeValidationError)
		return
	}

	if len(req.Bookings) == 0 {
		http_helpers.SendError(w, http.StatusBadRequest, "booking is empty", http_helpers.ErrorTypeValidationError)
		return
	}

	order := domain.Order{
		ID:     req.ID,
		UserID: req.UserID,
	}

	for _, book := range req.Bookings {
		if book.To.Before(book.From.Time) {
			http_helpers.SendError(w, http.StatusBadRequest,
				fmt.Sprintf("invalid date range for hotel id %v", book.HotelID),
				http_helpers.ErrorTypeValidationError)
			return
		}

		if !domain.RoomTypes.Contains(book.RoomType) {
			http_helpers.SendError(w, http.StatusBadRequest,
				fmt.Sprintf("invalid room_type '%s' for hotel id %v", book.RoomType, book.HotelID),
				http_helpers.ErrorTypeValidationError)
			return
		}

		newBooking := domain.Booking{
			HotelID:   book.HotelID,
			RoomType:  book.RoomType,
			From:      book.From.Time,
			To:        book.To.Time,
			RoomCount: book.RoomCount,
		}

		order.Bookings = append(order.Bookings, newBooking)
	}

	createdOrder, err := h.booking.CreateOrder(ctx, order)
	if err != nil {
		if errors.Is(err, domain.ErrOrderAlreadyExists) {
			http_helpers.SendSuccess(w, http.StatusOK, createdOrder)
			return
		}

		if errors.Is(err, domain.ErrHotelNotFound) {
			http_helpers.SendError(w, http.StatusBadRequest, "invalid hotel id", http_helpers.ErrorTypeValidationError)
			return
		}

		if errors.Is(err, domain.ErrRoomTypeNotFound) {
			http_helpers.SendError(w, http.StatusBadRequest, "invalid room type", http_helpers.ErrorTypeValidationError)
			return
		}

		if errors.Is(err, domain.ErrRoomsNotAvailable) {
			http_helpers.SendError(w, http.StatusBadRequest, err.Error(), http_helpers.ErrorTypeValidationError)
			return
		}

		log.Error("failed to create order", err)
		http_helpers.SendError(w, http.StatusInternalServerError, "failed to create order", http_helpers.ErrorTypeInternalError)
		return
	}

	http_helpers.SendSuccess(w, http.StatusCreated, createdOrder)

	log.WithField("order", createdOrder).Info("order successfully created")
}
