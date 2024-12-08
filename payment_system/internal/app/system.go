package app

import (
	"payment_system/pkg/models"
)

type System interface {
	Create(*models.PaymentRequest) (*models.PaymentResponse, error)
	Send(*models.PaymentWebhookRequest) error
}
