package api

import (
	"git.solsynth.dev/hydrogen/messaging/pkg/internal/models"
	"git.solsynth.dev/hydrogen/messaging/pkg/internal/services"
	"github.com/gofiber/contrib/websocket"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog/log"
)

func messageGateway(c *websocket.Conn) {
	user := c.Locals("user").(models.Account)

	// Push connection
	services.ClientRegister(user, c)
	log.Debug().Uint("user", user.ID).Msg("New websocket connection established...")

	// Event loop
	var task models.UnifiedCommand

	var messageType int
	var packet []byte
	var err error

	for {
		if messageType, packet, err = c.ReadMessage(); err != nil {
			break
		} else if err := jsoniter.Unmarshal(packet, &task); err != nil {
			_ = c.WriteMessage(messageType, models.UnifiedCommand{
				Action:  "error",
				Message: "unable to unmarshal your command, requires json request",
			}.Marshal())
			continue
		}

		message := services.DealCommand(task, user)

		if message != nil {
			if err = c.WriteMessage(messageType, message.Marshal()); err != nil {
				break
			}
		}
	}

	// Pop connection
	services.ClientUnregister(user, c)
	log.Debug().Uint("user", user.ID).Msg("A websocket connection disconnected...")
}
