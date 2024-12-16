package api

import (
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"payment_system/pkg/api/v1"
	"payment_system/pkg/models"
	"time"
)

func (s *PaymentServer) PostApiV1Payments(ctx echo.Context) error {
	var paymentRequest api.PaymentRequest
	if err := ctx.Bind(&paymentRequest); err != nil {
		log.Printf("Error binding request: %v\n", err)
		return ctx.JSON(http.StatusBadRequest, err)
	}
	request := fromApiPaymentRequest(&paymentRequest)
	response, err := s.paymentService.Create(request)
	if err != nil {
		log.Printf("Error sending request: %v\n", err)
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	WebhookRequest := &models.PaymentWebhookRequest{
		BookingId: request.BookingId,
		PaymentId: response.PaymentID,
		Status:    "ok",
	}
	go func() {
		time.Sleep(2 * time.Second)
		err = s.paymentService.Send(WebhookRequest)
		if err != nil {
			log.Printf("Error sending webhook: %v\n", err)
		}
	}()

	return ctx.JSON(http.StatusOK, response)
}
