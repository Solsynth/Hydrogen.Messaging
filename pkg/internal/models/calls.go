package models

import (
	"github.com/livekit/protocol/livekit"
	"time"
)

type Call struct {
	BaseModel

	EndedAt *time.Time `json:"ended_at"`

	ExternalID string        `json:"external_id"`
	FounderID  uint          `json:"founder_id"`
	ChannelID  uint          `json:"channel_id"`
	Founder    ChannelMember `json:"founder"`
	Channel    Channel       `json:"channel"`

	Participants []*livekit.ParticipantInfo `json:"participants" gorm:"-"`
}
