package api

import (
	"context"
	"github.com/ichinosekei/hotel-service-go/hotelier-service/internal/app"
	"github.com/ichinosekei/hotel-service-go/hotelier-service/proto"
	"log"
)

type HotelierServer struct {
	proto.UnimplementedHotelierServiceServer
	service app.HotelService
}

func NewHotelierServer(service app.HotelService) *HotelierServer {
	return &HotelierServer{service: service}
}

func (s *HotelierServer) GetRoomPrice(ctx context.Context, req *proto.RoomRequest) (*proto.RoomResponse, error) {
	roomPrice, err := s.service.GetRoomPrice(req.HotelId, req.RoomNumber)
	if err != nil {
		log.Printf("Error retrieving room price: %v", err)
		return nil, err
	}
	return &proto.RoomResponse{Price: roomPrice}, nil
}
