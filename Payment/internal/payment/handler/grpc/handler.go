package grpc_hadler

import (
	"context"

	pb "github.com/sewaustav/CaseGogRPServer/gen/go/payments"
	"github.com/sewaustav/Payment/internal/payment/models"
	service "github.com/sewaustav/Payment/internal/payment/service/api"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PaymentGRPCHandler struct {
	pb.UnimplementedPaymentCheckServiceServer
	service service.PaymentApiService
}

func NewPaymentGRPCHandler(service *service.PaymentApiCore) *PaymentGRPCHandler {
	return &PaymentGRPCHandler{
		service: service,
	}
}

func (h *PaymentGRPCHandler) CheckStatus(ctx context.Context, req *pb.UserInfo) (*pb.Status, error) {
	status, err := h.service.GetStatusService(ctx, models.UserIdentity{UserID: req.UserId})
	if err != nil {
		return nil, err
	}

	return &pb.Status{
		Status: int32(status.Status),
		ExpiredAt: timestamppb.New(status.ExpiredAt),
	}, nil
}