package http_handler

import (
	"github.com/gin-gonic/gin"
	service "github.com/sewaustav/Payment/internal/payment/service/api"
)

type PaymentHttpHandler struct {
	service service.PaymentApiService
}

func NewHttpHandler(service service.PaymentApiService) *PaymentHttpHandler {
	return &PaymentHttpHandler{
		service: service,
	}
}

