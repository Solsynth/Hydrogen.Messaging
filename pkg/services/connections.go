package services

import (
	"sync"

	"git.solsynth.dev/hydrogen/messaging/pkg/models"
	"github.com/gofiber/contrib/websocket"
)

var (
	wsMutex sync.Mutex
	wsConn  = make(map[uint]map[*websocket.Conn]bool)
)

func ClientRegister(user models.Account, conn *websocket.Conn) {
	wsMutex.Lock()
	if wsConn[user.ID] == nil {
		wsConn[user.ID] = make(map[*websocket.Conn]bool)
	}
	wsConn[user.ID][conn] = true
	wsMutex.Unlock()
}

func ClientUnregister(user models.Account, conn *websocket.Conn) {
	wsMutex.Lock()
	if wsConn[user.ID] == nil {
		wsConn[user.ID] = make(map[*websocket.Conn]bool)
	}
	delete(wsConn[user.ID], conn)
	wsMutex.Unlock()
}

func PushCommand(userId uint, task models.UnifiedCommand) {
	for conn := range wsConn[userId] {
		_ = conn.WriteMessage(1, task.Marshal())
	}
}

func DealCommand(task models.UnifiedCommand, user models.Account) *models.UnifiedCommand {
	switch task.Action {
	default:
		return &models.UnifiedCommand{
			Action:  "error",
			Message: "command not found",
		}
	}
}
