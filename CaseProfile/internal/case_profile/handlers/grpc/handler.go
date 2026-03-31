package grpc

import (
	"context"

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
		return nil, status.Error(codes.Internal, "user id not found in context")
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
		return nil, err
	}

	return &pb.Response{
		Status: "success",
	}, nil

}
