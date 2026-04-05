package grpc

import (
	"context"
	"errors"
	"testing"

	"github.com/sewaustav/CaseGoProfile/internal/case_profile/dto"
	"github.com/sewaustav/CaseGoProfile/internal/case_profile/models"
	"github.com/sewaustav/CaseGoProfile/mocks"
	"github.com/sewaustav/CaseGoProfile/pkg/middleware/rs256"
	pb "github.com/sewaustav/CaseGogRPServer/gen/go/case_go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func contextWithUser(userID int64, role int) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, rs256.UserIDKey, userID)
	ctx = context.WithValue(ctx, rs256.RoleKey, role)
	return ctx
}

func TestSendResult_Success(t *testing.T) {
	svc := mocks.NewService(t)
	handler := NewCaseGRPCHandler(svc)

	ctx := contextWithUser(42, int(models.User))

	ts := timestamppb.Now()
	req := &pb.CaseResult{
		CaseId:               1,
		DialogId:             10,
		StepsCount:           5,
		TokenUsed:            100,
		FinishedAt:           ts,
		Assertiveness:        0.8,
		Empathy:              0.7,
		ClarityCommunication: 0.6,
		Resistance:           0.5,
		Eloquence:            0.4,
		Initiative:           0.3,
	}

	expectedResult := dto.Result{
		UserID:               42,
		CaseID:               1,
		DialogID:             10,
		StepsCount:           5,
		TokensUsed:           100,
		FinishedAt:           ts.AsTime(),
		Assertiveness:        0.8,
		Empathy:              0.7,
		ClarityCommunication: 0.6,
		Resistance:           0.5,
		Eloquence:            0.4,
		Initiative:           0.3,
	}
	expectedUser := models.UserIdentity{UserID: 42, Role: models.User}

	svc.On("HandleResultsService", mock.Anything, expectedResult, expectedUser).Return(nil)

	resp, err := handler.SendResult(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "success", resp.Status)
}

func TestSendResult_NoUserID(t *testing.T) {
	svc := mocks.NewService(t)
	handler := NewCaseGRPCHandler(svc)

	ctx := context.Background() // нет user ID

	ts := timestamppb.Now()
	req := &pb.CaseResult{
		CaseId:     1,
		DialogId:   10,
		FinishedAt: ts,
	}

	resp, err := handler.SendResult(ctx, req)

	assert.Nil(t, resp)
	assert.Error(t, err)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Contains(t, st.Message(), "user id not found")
}

func TestSendResult_ServiceError(t *testing.T) {
	svc := mocks.NewService(t)
	handler := NewCaseGRPCHandler(svc)

	ctx := contextWithUser(42, int(models.User))

	ts := timestamppb.Now()
	req := &pb.CaseResult{
		CaseId:     1,
		DialogId:   10,
		FinishedAt: ts,
	}

	svc.On("HandleResultsService", mock.Anything, mock.AnythingOfType("dto.Result"), mock.AnythingOfType("models.UserIdentity")).
		Return(errors.New("service error"))

	resp, err := handler.SendResult(ctx, req)

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "service error")
}
