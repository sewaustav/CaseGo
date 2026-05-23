package grpc_handler

import (
	"context"

	dto "github.com/YoungFlores/Case_Go/Profile/internal/profile/dto"
	"github.com/YoungFlores/Case_Go/Profile/internal/profile/models"
	service "github.com/YoungFlores/Case_Go/Profile/internal/profile/service"
	pb "github.com/sewaustav/CaseGogRPServer/gen/go/level"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ProfileGrpcHandler struct {
	pb.UnimplementedLevelServer
	service service.ProfileCore
}

func NewProfileGrpcHandler(service service.ProfileCore) *ProfileGrpcHandler {
	return &ProfileGrpcHandler{
		service: service,
	}
}

func (h *ProfileGrpcHandler) LevelDone(ctx context.Context, req *pb.LevelResult) (*pb.UserInfo, error) {
	userProfile, err := h.service.UpdateLevelService(ctx, models.UserIdentity{
		UserID: req.UserId,
		Role:   models.User,
	}, &dto.LevelDto{
		Xp:   int(req.Xp),
		Date: req.DatePassed.AsTime(),
	})

	if err != nil {
		return nil, err
	}

	return &pb.UserInfo{
		UserId:     userProfile.UserID,
		Xp:         int32(userProfile.Xp),
		Level:      int32(userProfile.Level),
		Streak:     int32(userProfile.Streak),
		LastActive: timestamppb.New(userProfile.LastActive),
	}, nil
}
