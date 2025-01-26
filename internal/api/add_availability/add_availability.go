package add_availability

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"applicationDesignTest/internal/api/http_helpers"
	"applicationDesignTest/internal/domain"
	"applicationDesignTest/pkg/date"
	"applicationDesignTest/pkg/log"
)

type request struct {
	HotelID   domain.HotelID  `json:"hotel_id"`
	RoomType  domain.RoomType `json:"room_type"`
	Date      date.CustomDate `json:"date"`
	RoomCount int             `json:"room_count"`
}

type BookingService interface {
	AddRoomAvailability(ctx context.Context, hotelID domain.HotelID, roomType domain.RoomType, date time.Time, rooms int) error
}

type Handler struct {
	booking BookingService
}

func NewHandler(bs BookingService) *Handler {
	return &Handler{
		booking: bs,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req request

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http_helpers.SendError(w, http.StatusBadRequest, "invalid input", http_helpers.ErrorTypeValidationError)
		return
	}

	if err := h.booking.AddRoomAvailability(ctx, req.HotelID, req.RoomType, req.Date.Time, req.RoomCount); err != nil {
		if errors.Is(err, domain.ErrHotelNotFound) {
			http_helpers.SendError(w, http.StatusBadRequest, "invalid hotel id", http_helpers.ErrorTypeValidationError)
			return
		}

		log.Error("failed to add availability", err)
		http_helpers.SendError(w, http.StatusInternalServerError, err.Error(), http_helpers.ErrorTypeInternalError)
		return
	}

	http_helpers.SendSuccess(w, http.StatusOK, nil)
}
