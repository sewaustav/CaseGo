package grpc_handler

import (
	pb "github.com/sewaustav/CaseGogRPServer/gen/go/level"
)

type ProfileGrpcHandler struct {
	pb.UnimplementedProfileServer
}
