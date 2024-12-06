package models

type PaymentRequest struct {
	Amount    float64
	BookingId string
}

type PaymentResponse struct {
	PaymentId string
}
