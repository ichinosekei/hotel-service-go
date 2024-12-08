package api

import (
	"payment_system/pkg/api/v1"
	"payment_system/pkg/models"
)

func fromApiPaymentRequest(paymentRequest *api.PaymentRequest) *models.PaymentRequest {
	return &models.PaymentRequest{
		Amount:    paymentRequest.Amount,
		BookingId: paymentRequest.BookingId,
	}
}
