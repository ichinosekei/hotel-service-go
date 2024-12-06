package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type PaymentRequest struct {
	Amount    float64 `json:"amount"`
	BookingId string  `json:"bookingId"`
}

type PaymentClient struct {
	req  PaymentRequest
	addr string
}

func NewPaymentClient(addr string, amount float64, bookingId string) *PaymentClient {
	return &PaymentClient{
		req: PaymentRequest{
			Amount:    amount,
			BookingId: bookingId,
		},
		addr: addr,
	}
}
func (p *PaymentClient) InitiatePayment() error {
	data, err := json.Marshal(p.req)
	if err != nil {
		log.Printf("Failed to marshalling payment request: %v", err)
		return err
	}
	resp, err := http.Post(p.addr, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Printf("Failde to posting payment request: %v", err)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Payment request failed with status %v", resp.StatusCode)
		return fmt.Errorf("payment system responded with status: %v", resp.Status)
	}
	return nil
}
