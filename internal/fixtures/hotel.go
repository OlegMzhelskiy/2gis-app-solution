package fixtures

import (
	"context"
	"time"

	"applicationDesignTest/internal/domain"
)

type hotelRepository interface {
	AddHotel(ctx context.Context, hotel domain.Hotel) error
	AddRoomAvailability(ctx context.Context, hotelID domain.HotelID, roomType domain.RoomType, date time.Time, rooms int) error
}

func InitHotelData(store hotelRepository) error {
	ctx := context.Background()

	var err error

	reddison := domain.Hotel{
		ID:   1,
		Name: "Reddison",
	}

	if err = store.AddHotel(ctx, reddison); err != nil {
		return err
	}

	if err = store.AddRoomAvailability(ctx, reddison.ID, domain.RoomTypeSingle,
		time.Date(2025, time.February, 1, 0, 0, 0, 0, time.UTC), 1); err != nil {
		return err
	}
	if err = store.AddRoomAvailability(ctx, reddison.ID, domain.RoomTypeSingle,
		time.Date(2025, time.February, 2, 0, 0, 0, 0, time.UTC), 3); err != nil {
		return err
	}
	if err = store.AddRoomAvailability(ctx, reddison.ID, domain.RoomTypeSingle,
		time.Date(2025, time.February, 3, 0, 0, 0, 0, time.UTC), 3); err != nil {
		return err
	}
	if err = store.AddRoomAvailability(ctx, reddison.ID, domain.RoomTypeSingle,
		time.Date(2025, time.February, 4, 0, 0, 0, 0, time.UTC), 3); err != nil {
		return err
	}
	if err = store.AddRoomAvailability(ctx, reddison.ID, domain.RoomTypeSingle,
		time.Date(2025, time.February, 5, 0, 0, 0, 0, time.UTC), 3); err != nil {
		return err
	}

	return nil
}
