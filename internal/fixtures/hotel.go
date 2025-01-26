package fixtures

import (
	"time"

	"applicationDesignTest/internal/domain"
)

type hotelRepository interface {
	AddHotel(hotel domain.Hotel) error
	AddRoomAvailability(hotelID domain.HotelID, roomType domain.RoomType, date time.Time, rooms int) error
}

func InitHotelData(store hotelRepository) error {
	var err error

	reddison := domain.Hotel{
		ID:   1,
		Name: "Reddison",
	}

	if err = store.AddHotel(reddison); err != nil {
		return err
	}

	if err = store.AddRoomAvailability(reddison.ID, domain.RoomTypeSingle,
		time.Date(2025, time.February, 1, 0, 0, 0, 0, time.UTC), 1); err != nil {
		return err
	}
	if err = store.AddRoomAvailability(reddison.ID, domain.RoomTypeSingle,
		time.Date(2025, time.February, 2, 0, 0, 0, 0, time.UTC), 3); err != nil {
		return err
	}
	if err = store.AddRoomAvailability(reddison.ID, domain.RoomTypeSingle,
		time.Date(2025, time.February, 3, 0, 0, 0, 0, time.UTC), 3); err != nil {
		return err
	}
	if err = store.AddRoomAvailability(reddison.ID, domain.RoomTypeSingle,
		time.Date(2025, time.February, 4, 0, 0, 0, 0, time.UTC), 3); err != nil {
		return err
	}
	if err = store.AddRoomAvailability(reddison.ID, domain.RoomTypeSingle,
		time.Date(2025, time.February, 5, 0, 0, 0, 0, time.UTC), 3); err != nil {
		return err
	}

	return nil
}
