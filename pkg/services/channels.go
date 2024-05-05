package services

import (
	"fmt"
	"regexp"

	"git.solsynth.dev/hydrogen/messaging/pkg/database"
	"git.solsynth.dev/hydrogen/messaging/pkg/models"
	"github.com/samber/lo"
)

func GetChannelAliasAvailability(alias string) error {
	if !regexp.MustCompile("^[a-z0-9-]+$").MatchString(alias) {
		return fmt.Errorf("channel alias should only contains lowercase letters, numbers, and hyphens")
	}
	return nil
}

func GetChannel(id uint) (models.Channel, error) {
	var channel models.Channel
	if err := database.C.Where(models.Channel{
		BaseModel: models.BaseModel{ID: id},
	}).Preload("Account").First(&channel).Error; err != nil {
		return channel, err
	}

	return channel, nil
}

func GetChannelWithAlias(alias string, realmId ...uint) (models.Channel, error) {
	var channel models.Channel
	tx := database.C.Where(models.Channel{Alias: alias}).Preload("Account")
	if len(realmId) > 0 {
		tx = tx.Where("realm_id = ?", realmId)
	} else {
		tx = tx.Where("realm_id IS NULL")
	}
	if err := tx.First(&channel).Error; err != nil {
		return channel, err
	}

	return channel, nil
}

func GetAvailableChannelWithAlias(alias string, user models.Account, realmId ...uint) (models.Channel, models.ChannelMember, error) {
	var err error
	var member models.ChannelMember
	var channel models.Channel
	if channel, err = GetChannelWithAlias(alias, realmId...); err != nil {
		return channel, member, err
	}

	if err := database.C.Where(models.ChannelMember{
		AccountID: user.ID,
		ChannelID: channel.ID,
	}).First(&member).Error; err != nil {
		return channel, member, fmt.Errorf("channel principal not found: %v", err.Error())
	}

	return channel, member, nil
}

func GetAvailableChannel(id uint, user models.Account) (models.Channel, models.ChannelMember, error) {
	var err error
	var member models.ChannelMember
	var channel models.Channel
	if channel, err = GetChannel(id); err != nil {
		return channel, member, err
	}

	if err := database.C.Where(models.ChannelMember{
		AccountID: user.ID,
		ChannelID: channel.ID,
	}).First(&member).Error; err != nil {
		return channel, member, fmt.Errorf("channel principal not found: %v", err.Error())
	}

	return channel, member, nil
}

func ListChannel(realmId ...uint) ([]models.Channel, error) {
	var channels []models.Channel
	tx := database.C.Preload("Account")
	if len(realmId) > 0 {
		tx = tx.Where("realm_id = ?", realmId)
	} else {
		tx = tx.Where("realm_id IS NULL")
	}
	if err := tx.Find(&channels).Error; err != nil {
		return channels, err
	}

	return channels, nil
}

func ListChannelWithUser(user models.Account, realmId ...uint) ([]models.Channel, error) {
	var channels []models.Channel
	tx := database.C.Where(&models.Channel{AccountID: user.ID})
	if len(realmId) > 0 {
		tx = tx.Where("realm_id = ?", realmId)
	} else {
		tx = tx.Where("realm_id IS NULL")
	}
	if err := tx.Find(&channels).Error; err != nil {
		return channels, err
	}

	return channels, nil
}

func ListAvailableChannel(user models.Account, realmId ...uint) ([]models.Channel, error) {
	var channels []models.Channel
	var members []models.ChannelMember
	if err := database.C.Where(&models.ChannelMember{
		AccountID: user.ID,
	}).Find(&members).Error; err != nil {
		return channels, err
	}

	idx := lo.Map(members, func(item models.ChannelMember, index int) uint {
		return item.ChannelID
	})

	tx := database.C.Where("id IN ?", idx)
	if len(realmId) > 0 {
		tx = tx.Where("realm_id = ?", realmId)
	} else {
		tx = tx.Where("realm_id IS NULL")
	}
	if err := tx.Find(&channels).Error; err != nil {
		return channels, err
	}

	return channels, nil
}

func NewChannel(user models.Account, alias, name, description string, realmId ...uint) (models.Channel, error) {
	channel := models.Channel{
		Alias:       alias,
		Name:        name,
		Description: description,
		AccountID:   user.ID,
		Members: []models.ChannelMember{
			{AccountID: user.ID},
		},
	}
	if len(realmId) > 0 {
		channel.RealmID = &realmId[0]
	}

	err := database.C.Save(&channel).Error

	return channel, err
}

func EditChannel(channel models.Channel, alias, name, description string) (models.Channel, error) {
	channel.Alias = alias
	channel.Name = name
	channel.Description = description

	err := database.C.Save(&channel).Error

	return channel, err
}

func DeleteChannel(channel models.Channel) error {
	return database.C.Delete(&channel).Error
}
