package hub

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"
)

type WebSse interface {
	Register(userID int64, role string) Client
	UnRegister(client Client)
	Broadcast(msg string)
	SendMessage(client Client) string
}

type webSse struct {
}

var webSseInstance WebSse
var webSseOnce sync.Once

func NewWebSse() WebSse {
	webSseOnce.Do(func() {
		webSseInstance = webSse{}
	})
	return webSseInstance
}

func (w webSse) Register(userID int64, role string) Client {
	client := Client{
		UUID:    fmt.Sprintf("%d", time.Now().UnixNano()),
		UserID:  userID,
		Role:    role,
		Message: make(chan Message, 100),
	}
	hub.Register <- client
	return client
}

func (w webSse) UnRegister(client Client) {
	hub.UnRegister <- client
}

func (w webSse) Broadcast(msg string) {
	var hubMsg Message
	if err := json.Unmarshal([]byte(msg), &hubMsg); err != nil {
		log.Println("webSse Broadcast Unmarshal error", err)
		return
	}
	hub.Broadcast <- Message{
		Topic:   hubMsg.Topic,
		Role:    hubMsg.Role,
		SendID:  hubMsg.SendID,
		UserID:  hubMsg.UserID,
		Content: hubMsg.Content,
	}
}

func (w webSse) SendMessage(client Client) string {
	select {
	case message := <-client.Message:
		return message.Content
	default:
		log.Println(fmt.Sprintf("No message for client %d : %s", client.UserID, client.UUID))
		return ""
	}
}
