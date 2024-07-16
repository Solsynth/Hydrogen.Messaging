package services

import (
	"context"
	"git.solsynth.dev/hydrogen/dealer/pkg/proto"
	"git.solsynth.dev/hydrogen/messaging/pkg/internal/gap"
	"git.solsynth.dev/hydrogen/messaging/pkg/internal/models"
	"time"
)

func PushCommand(userId uint, task models.UnifiedCommand) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pc := gap.H.GetDealerGrpcConn()
	_, _ = proto.NewStreamControllerClient(pc).PushStream(ctx, &proto.PushStreamRequest{
		UserId: uint64(userId),
		Body:   task.Marshal(),
	})
}
