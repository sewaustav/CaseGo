package grpc

import (
	"context"

	pb "github.com/sewaustav/CaseGogRPServer/gen/go/case_go"
)

type CaseGRPCHandler struct {
	pb.UnimplementedCasesServer
}

func NewCaseGRPCHandler() *CaseGRPCHandler {
	return &CaseGRPCHandler{}
}

func (h *CaseGRPCHandler) SendResult(ctx context.Context, req *pb.CaseResult) (*pb.Response, error) {

	return nil, nil
}
