package domain

import "errors"

var (
	ErrHotelNotFound      = errors.New("hotel not found")
	ErrRoomTypeNotFound   = errors.New("room type not found")
	ErrOrderNotFound      = errors.New("order not found")
	ErrOrderAlreadyExists = errors.New("order already exists")
	ErrRoomsNotAvailable  = errors.New("rooms not available")
)
