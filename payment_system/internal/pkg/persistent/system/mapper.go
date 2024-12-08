package system

import (
	"payment_system/pkg/models"
)

func toModelPaymentResponse(resp *PaymentResponse) *models.PaymentResponse {
	return &models.PaymentResponse{
		URL:       resp.URL,
		PaymentID: resp.PaymentID,
	}
}

func fromModelPaymentRequest(req *models.PaymentRequest) *PaymentRequest {
	return &PaymentRequest{
		Amount:    req.Amount,
		BookingId: req.BookingId,
	}
}

func fromModelPaymentWebhookRequest(req *models.PaymentWebhookRequest) *PaymentWebhookRequest {
	return &PaymentWebhookRequest{
		BookingId: req.BookingId,
		PaymentId: req.PaymentId,
		Status:    req.Status,
	}
}
