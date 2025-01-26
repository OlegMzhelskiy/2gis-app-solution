package domain

import "time"

type OrderNumber int64

type OrderID string

type Order struct {
	ID        OrderID     `json:"id"`
	Number    OrderNumber `json:"number"`
	UserID    UserID      `json:"user_id"`
	CreatedAt time.Time   `json:"created_at"`
	Bookings  []Booking   `json:"booking"`
}

type Booking struct {
	HotelID   HotelID   `json:"hotel_id"`
	RoomType  RoomType  `json:"room_type"`
	From      time.Time `json:"from"`
	To        time.Time `json:"to"`
	RoomCount int       `json:"room_count"`
}
