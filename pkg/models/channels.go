package models

type ChannelType = uint8

const (
	ChannelTypeDirect = ChannelType(iota)
	ChannelTypeRealm
)

type Channel struct {
	BaseModel

	Alias       string          `json:"alias" gorm:"uniqueIndex"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Members     []ChannelMember `json:"members"`
	Messages    []Message       `json:"messages"`
	Type        ChannelType     `json:"type"`
	AccountID   uint            `json:"account_id"`
	RealmID     uint            `json:"realm_id"`
}

type ChannelMember struct {
	BaseModel

	ChannelID uint    `json:"channel_id"`
	AccountID uint    `json:"account_id"`
	Channel   Channel `json:"channel"`
	Account   Account `json:"account"`

	Messages []Message `json:"messages" gorm:"foreignKey:SenderID"`
}
