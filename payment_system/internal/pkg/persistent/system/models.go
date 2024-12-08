package system

type PaymentRequest struct {
	Amount    float64
	BookingId string
}

type PaymentResponse struct {
	URL       string
	PaymentID string
}

type PaymentWebhookRequest struct {
	BookingId string `json:"bookingId"`
	PaymentId string `json:"paymentId"`
	Status    string `json:"status"`
}
