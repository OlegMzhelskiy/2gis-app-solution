package main

import (
	"fmt"
	"net/http"

	"applicationDesignTest/internal/api/add_availability"
	"applicationDesignTest/internal/api/create_order"
	"applicationDesignTest/internal/api/get_order"
	"applicationDesignTest/internal/config"
	"applicationDesignTest/internal/fixtures"
	"applicationDesignTest/internal/storage/memorystore"
	"applicationDesignTest/internal/usecase/booking"
	"applicationDesignTest/internal/usecase/order"
	"applicationDesignTest/pkg/log"

	"github.com/go-chi/chi/v5"

	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	log.InitializeLogger()

	l := log.GetLogger()
	defer l.Sync() //nolint:errcheck

	if err := run(); err != nil {
		log.Error("server stopped", err)
	}
}

func run() error {
	log.Info("init config")

	cfg, err := config.LoadConfig(".")
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	log.Info("init store")

	hotelStore := memorystore.NewHotelStore()
	orderStore := memorystore.NewOrderStore()

	orderService := order.NewOrderService(orderStore)
	bookingService := booking.NewBookingService(hotelStore, orderService)

	getOrderHandler := get_order.NewHandler(orderStore)
	createOrderHandler := create_order.NewHandler(bookingService)
	addAvailabilityHandler := add_availability.NewHandler(bookingService)

	log.Info("init fixtures")

	if err := fixtures.InitHotelData(hotelStore); err != nil {
		return fmt.Errorf("can't init fixtures: %w", err)
	}

	log.Info("register handlers")

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/orders/{orderNumber}", getOrderHandler.Handle)
	r.Post("/orders", createOrderHandler.Handle)
	r.Post("/hotels/availability", addAvailabilityHandler.Handle)

	log.Info(fmt.Sprintf("server is running on port %v", cfg.Port))

	if err = http.ListenAndServe(":"+cfg.Port, r); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}
