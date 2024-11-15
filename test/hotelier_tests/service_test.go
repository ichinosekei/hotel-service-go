package hotelier_tests

import (
	"hotel-service-go/internal/hotelier"
	"testing"
)

func TestCreateHotel(t *testing.T) {
	service := hotelier.NewService(setupTestDB())
	hotelID, err := service.CreateHotel("Test Hotel", "City")
	if err != nil || hotelID == 0 {
		t.Fatalf("expected to create hotel, got error: %v", err)
	}
}
