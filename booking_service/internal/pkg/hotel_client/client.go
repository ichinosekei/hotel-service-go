package hotel_client

import (
	"context"
	"github.com/ichinosekei/hotel-service-go/hotelier-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"strconv"
	"time"
)

type HotelClient struct {
	client proto.HotelierServiceClient
	conn   *grpc.ClientConn
}

func NewHotelClient(address string) (*HotelClient, error) {
	connection, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect to HotelService: %v", err)
		return nil, err
	}
	return &HotelClient{
		client: proto.NewHotelierServiceClient(connection),
		conn:   connection}, nil
}

func (hc *HotelClient) GetRoomPrice(hotelID int, roomNumber int) (float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := hc.client.GetRoomPrice(ctx, &proto.RoomRequest{
		HotelId:    int32(hotelID),
		RoomNumber: strconv.Itoa(roomNumber),
	})
	if err != nil {
		log.Printf("Failed to get room price: %v", err)
		return 0, err
	}
	return resp.Price, nil
}

func (hc *HotelClient) Close() error {
	return hc.conn.Close()
}
