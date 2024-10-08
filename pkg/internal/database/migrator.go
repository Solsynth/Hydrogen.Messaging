package database

import (
	"git.solsynth.dev/hydrogen/messaging/pkg/internal/models"
	"gorm.io/gorm"
)

var AutoMaintainRange = []any{
	&models.Account{},
	&models.Realm{},
	&models.Channel{},
	&models.ChannelMember{},
	&models.Call{},
	&models.Event{},
}

func RunMigration(source *gorm.DB) error {
	if err := source.AutoMigrate(AutoMaintainRange...); err != nil {
		return err
	}

	return nil
}
