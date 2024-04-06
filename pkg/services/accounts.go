package services

import (
	"context"
	"time"

	"git.solsynth.dev/hydrogen/identity/pkg/grpc/proto"
	"git.solsynth.dev/hydrogen/messaging/pkg/grpc"
	"git.solsynth.dev/hydrogen/messaging/pkg/models"
	"github.com/spf13/viper"
)

func GetAccountFriend(userId, relatedId uint, status int) (*proto.FriendshipResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	return grpc.Friendships.GetFriendship(ctx, &proto.FriendshipTwoSideLookupRequest{
		AccountId: uint64(userId),
		RelatedId: uint64(relatedId),
		Status:    uint32(status),
	})
}

func NotifyAccount(user models.Account, subject, content string, realtime bool, links ...*proto.NotifyLink) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := grpc.Notify.NotifyUser(ctx, &proto.NotifyRequest{
		ClientId:     viper.GetString("identity.client_id"),
		ClientSecret: viper.GetString("identity.client_secret"),
		Subject:      subject,
		Content:      content,
		Links:        links,
		RecipientId:  uint64(user.ExternalID),
		IsRealtime:   realtime,
		IsImportant:  false,
	})

	return err
}
