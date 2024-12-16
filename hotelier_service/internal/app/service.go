package app

type HotelService interface {
	GetRoomPrice(hotelID int32, roomNumber string) (float64, error)
}
