package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	defer func() {
		if err := l.Sync(); err != nil {
			log.Error("failed to sync logger", err)
		}
	}()

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

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("server failed", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	log.Info("received shutdown signal, shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("server shutdown failed", err)
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	log.Info("server gracefully stopped")

	return nil
}
