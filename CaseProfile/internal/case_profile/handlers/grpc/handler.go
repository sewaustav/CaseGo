package grpc

import (
	"context"
	"errors"
	"net/http"

	"github.com/sewaustav/CaseGoProfile/apperrors"
	"github.com/sewaustav/CaseGoProfile/internal/case_profile/dto"
	"github.com/sewaustav/CaseGoProfile/internal/case_profile/models"
	"github.com/sewaustav/CaseGoProfile/internal/case_profile/service"
	"github.com/sewaustav/CaseGoProfile/pkg/middleware/rs256"
	pb "github.com/sewaustav/CaseGogRPServer/gen/go/case_go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CaseGRPCHandler struct {
	pb.UnimplementedCasesServer
	service service.Service
}

func NewCaseGRPCHandler(service service.Service) *CaseGRPCHandler {
	return &CaseGRPCHandler{
		service: service,
	}
}

func (h *CaseGRPCHandler) SendResult(ctx context.Context, req *pb.CaseResult) (*pb.Response, error) {
	ts := req.FinishedAt
	t := ts.AsTime()

	userID, ok := ctx.Value(rs256.UserIDKey).(int64)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user id not found in context")
	}

	role, _ := ctx.Value(rs256.RoleKey).(int)

	result := &dto.Result{
		UserID:               userID,
		CaseID:               req.CaseId,
		DialogID:             req.DialogId,
		StepsCount:           req.StepsCount,
		TokensUsed:           req.TokenUsed,
		FinishedAt:           t,
		Assertiveness:        req.Assertiveness,
		Empathy:              req.Empathy,
		ClarityCommunication: req.ClarityCommunication,
		Resistance:           req.Resistance,
		Eloquence:            req.Eloquence,
		Initiative:           req.Initiative,
	}

	user := models.UserIdentity{
		UserID: userID,
		Role:   models.UserRole(role),
	}

	if err := h.service.HandleResultsService(ctx, *result, user); err != nil {
		return nil, mapAppErrorToGRPC(err)
	}

	return &pb.Response{
		Status: "success",
	}, nil
}

// mapAppErrorToGRPC конвертирует *apperrors.AppError в gRPC status.Error
func mapAppErrorToGRPC(err error) error {
	var appErr *apperrors.AppError
	if !errors.As(err, &appErr) {
		return status.Error(codes.Internal, "internal server error")
	}

	switch appErr.Code {
	case http.StatusNotFound:
		return status.Error(codes.NotFound, appErr.Message)
	case http.StatusForbidden:
		return status.Error(codes.PermissionDenied, appErr.Message)
	case http.StatusBadRequest:
		return status.Error(codes.InvalidArgument, appErr.Message)
	case http.StatusConflict:
		return status.Error(codes.AlreadyExists, appErr.Message)
	case http.StatusUnauthorized:
		return status.Error(codes.Unauthenticated, appErr.Message)
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
