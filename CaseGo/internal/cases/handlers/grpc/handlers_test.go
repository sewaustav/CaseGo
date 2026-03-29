package grpc

import (
	"context"
	"testing"
	"time"

	"github.com/sewaustav/CaseGoCore/internal/cases/models"
	pb "github.com/sewaustav/CaseGogRPServer/gen/go/case_go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

type fakeCasesClient struct {
	lastReq *pb.CaseResult
}

func (f *fakeCasesClient) SendResult(ctx context.Context, req *pb.CaseResult, opts ...grpc.CallOption) (*pb.Response, error) {
	f.lastReq = req
	return &pb.Response{}, nil
}

func TestSendResults(t *testing.T) {
	client := &fakeCasesClient{}
	h := &CaseGoGRPC{client: client}

	now := time.Now().UTC().Truncate(time.Second)
	err := h.SendResults(context.Background(), models.Result{
		UserID:               10,
		DialogID:             20,
		CaseID:               30,
		StepsCount:           4,
		TokensUsed:           100,
		FinishedAt:           now,
		Assertiveness:        1.1,
		Empathy:              2.2,
		ClarityCommunication: 3.3,
		Resistance:           4.4,
		Eloquence:            5.5,
		Initiative:           6.6,
	})

	require.NoError(t, err)
	require.NotNil(t, client.lastReq)
	assert.Equal(t, int64(10), client.lastReq.UserId)
	assert.Equal(t, int64(20), client.lastReq.DialogId)
	assert.Equal(t, int64(30), client.lastReq.CaseId)
	assert.Equal(t, int32(4), client.lastReq.StepsCount)
	assert.Equal(t, int32(100), client.lastReq.TokenUsed)
}
