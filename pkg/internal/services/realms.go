package services

import (
	"context"
	"errors"
	"fmt"
	"git.solsynth.dev/hydrogen/dealer/pkg/hyper"
	"git.solsynth.dev/hydrogen/dealer/pkg/proto"
	"git.solsynth.dev/hydrogen/messaging/pkg/internal/database"
	"git.solsynth.dev/hydrogen/messaging/pkg/internal/gap"
	"git.solsynth.dev/hydrogen/messaging/pkg/internal/models"
	"github.com/samber/lo"
	"gorm.io/gorm"
	"reflect"
)

func GetRealmWithExtID(id uint) (models.Realm, error) {
	var realm models.Realm
	pc, err := gap.H.GetServiceGrpcConn(hyper.ServiceTypeAuthProvider)
	if err != nil {
		return realm, err
	}
	response, err := proto.NewRealmClient(pc).GetRealm(context.Background(), &proto.LookupRealmRequest{
		Id: lo.ToPtr(uint64(id)),
	})
	if err != nil {
		return realm, err
	}
	return LinkRealm(response)
}

func GetRealmWithAlias(alias string) (models.Realm, error) {
	var realm models.Realm
	pc, err := gap.H.GetServiceGrpcConn(hyper.ServiceTypeAuthProvider)
	if err != nil {
		return realm, err
	}
	response, err := proto.NewRealmClient(pc).GetRealm(context.Background(), &proto.LookupRealmRequest{
		Alias: &alias,
	})
	if err != nil {
		return realm, err
	}
	return LinkRealm(response)
}

func GetRealmMember(realmId uint, userId uint) (*proto.RealmMemberInfo, error) {
	var realm models.Realm
	if err := database.C.Where("id = ?", realmId).First(&realm).Error; err != nil {
		return nil, err
	}
	pc, err := gap.H.GetServiceGrpcConn(hyper.ServiceTypeAuthProvider)
	if err != nil {
		return nil, err
	}
	response, err := proto.NewRealmClient(pc).GetRealmMember(context.Background(), &proto.RealmMemberLookupRequest{
		RealmId: uint64(realm.ExternalID),
		UserId:  lo.ToPtr(uint64(userId)),
	})
	if err != nil {
		return nil, err
	} else {
		return response, nil
	}
}

func ListRealmMember(realmId uint) ([]*proto.RealmMemberInfo, error) {
	pc, err := gap.H.GetServiceGrpcConn(hyper.ServiceTypeAuthProvider)
	if err != nil {
		return nil, err
	}
	response, err := proto.NewRealmClient(pc).ListRealmMember(context.Background(), &proto.RealmMemberLookupRequest{
		RealmId: uint64(realmId),
	})
	if err != nil {
		return nil, err
	} else {
		return response.Data, nil
	}
}

func LinkRealm(info *proto.RealmInfo) (models.Realm, error) {
	var realm models.Realm
	if info == nil {
		return realm, fmt.Errorf("remote realm info was not found")
	}
	if err := database.C.Where(&models.Realm{
		ExternalID: uint(info.Id),
	}).First(&realm).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			realm = models.Realm{
				Alias:       info.Alias,
				Name:        info.Name,
				Description: info.Description,
				IsPublic:    info.IsPublic,
				IsCommunity: info.IsCommunity,
				ExternalID:  uint(info.Id),
			}
			return realm, database.C.Save(&realm).Error
		}
		return realm, err
	}

	prev := realm
	realm.Alias = info.Alias
	realm.Name = info.Name
	realm.Description = info.Description
	realm.IsPublic = info.IsPublic
	realm.IsCommunity = info.IsCommunity

	var err error
	if !reflect.DeepEqual(prev, realm) {
		err = database.C.Save(&realm).Error
	}

	return realm, err
}
