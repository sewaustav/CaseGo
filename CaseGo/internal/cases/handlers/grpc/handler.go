package grpc

import (
	"context"

	"github.com/sewaustav/CaseGoCore/internal/cases/models"
	"github.com/sewaustav/CaseGoCore/internal/jwt"
	pb "github.com/sewaustav/CaseGogRPServer/gen/go/case_go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GRPCService interface {
	SendResults(ctx context.Context, msg models.Result) error
}

type CaseGoGRPC struct {
	token  jwt.JwtService
	client pb.CasesClient
}

func NewCaseGoGRPC(addr string, token jwt.JwtService) (*CaseGoGRPC, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &CaseGoGRPC{
		client: pb.NewCasesClient(conn),
		token:  token,
	}, nil
}

func (c *CaseGoGRPC) SendResults(ctx context.Context, msg models.Result) error {
	jwtToken, err := c.token.GenerateToken(msg.UserID, models.User)
	if err != nil {
		return err
	}

	md := metadata.Pairs("authorization", "Bearer"+jwtToken)
	authCtx := metadata.NewOutgoingContext(ctx, md)

	req := &pb.CaseResult{
		UserId:               msg.UserID,
		DialogId:             msg.DialogID,
		CaseId:               msg.CaseID,
		StepsCount:           msg.StepsCount,
		TokenUsed:            msg.TokensUsed,
		FinishedAt:           timestamppb.New(msg.FinishedAt),
		Assertiveness:        float32(msg.Assertiveness),
		Empathy:              float32(msg.Empathy),
		ClarityCommunication: float32(msg.ClarityCommunication),
		Resistance:           float32(msg.Resistance),
		Eloquence:            float32(msg.Eloquence),
		Initiative:           float32(msg.Initiative),
	}
	_, err = c.client.SendResult(authCtx, req)
	return err
}
