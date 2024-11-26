package grpc

import (
	"context"
	"time"

	"github.com/ichinosekei/hotel-service-go/hotelier-service/internal/pkg/proto"
	"google.golang.org/grpc"
)

type HotelierClient struct {
	client proto.HotelierServiceClient
}

func NewHotelierClient(address string) (*HotelierClient, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	return &HotelierClient{
		client: proto.NewHotelierServiceClient(conn),
	}, nil
}

func (hc *HotelierClient) GetRoomPrice(hotelID int32, roomNumber string) (float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := hc.client.GetRoomPrice(ctx, &proto.RoomRequest{
		HotelId:    hotelID,
		RoomNumber: roomNumber,
	})
	if err != nil {
		return 0, err
	}
	return resp.Price, nil
}
