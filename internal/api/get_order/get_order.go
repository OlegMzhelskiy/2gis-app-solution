package get_order

import (
	"errors"
	"net/http"
	"strconv"

	"applicationDesignTest/internal/api/http_helpers"
	"applicationDesignTest/internal/domain"
	"applicationDesignTest/pkg/log"
	"github.com/go-chi/chi/v5"
)

type OrderService interface {
	GetOrderByNumber(orderNumber domain.OrderNumber) (*domain.Order, error)
}

type Handler struct {
	orderService OrderService
}

func NewHandler(orderService OrderService) *Handler {
	return &Handler{
		orderService: orderService,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	orderNumberStr := chi.URLParam(r, "orderNumber")

	orderNumber, err := strconv.Atoi(orderNumberStr)
	if err != nil {
		http_helpers.SendError(w, http.StatusBadRequest, "invalid order number", http_helpers.ErrorTypeValidationError)
	}

	order, err := h.orderService.GetOrderByNumber(domain.OrderNumber(orderNumber))
	if err != nil {
		if errors.Is(err, domain.ErrOrderNotFound) {
			http_helpers.SendError(w, http.StatusBadRequest, "such order doesn't exist", http_helpers.ErrorTypeValidationError)
			return
		}

		log.Error("failed to add availability", err)
		http_helpers.SendError(w, http.StatusInternalServerError, err.Error(), http_helpers.ErrorTypeInternalError)
		return
	}

	http_helpers.SendSuccess(w, http.StatusOK, order)
}
