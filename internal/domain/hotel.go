package domain

type HotelID int

type Hotel struct {
	ID   HotelID
	Name string
}

type RoomType string

const (
	RoomTypeSingle RoomType = "single"
	RoomTypeDouble RoomType = "double"
	RoomTypeLux    RoomType = "lux"
)

type RoomTypesEnum map[RoomType]struct{}

var RoomTypes = RoomTypesEnum{
	RoomTypeSingle: {},
	RoomTypeDouble: {},
	RoomTypeLux:    {},
}

func (r RoomTypesEnum) Contains(value RoomType) bool {
	_, ok := r[value]
	return ok
}
