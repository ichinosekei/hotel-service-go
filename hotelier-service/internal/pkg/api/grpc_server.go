package api

import (
	"context"
	"github.com/ichinosekei/hotel-service-go/hotelier-service/internal/pkg/proto"
	"github.com/ichinosekei/hotel-service-go/hotelier-service/internal/pkg/repository"
	"log"
)

type HotelierServer struct {
	proto.UnimplementedHotelierServiceServer
	Service *repository.Service
}

func NewHotelierServer(service *repository.Service) *HotelierServer {
	return &HotelierServer{Service: service}
}

func (s *HotelierServer) GetRoomPrice(ctx context.Context, req *proto.RoomRequest) (*proto.RoomResponse, error) {
	roomPrice, err := s.Service.GetRoomPrice(req.HotelId, req.RoomNumber)
	if err != nil {
		log.Printf("Error retrieving room price: %v", err)
		return nil, err
	}
	return &proto.RoomResponse{Price: roomPrice}, nil
}
