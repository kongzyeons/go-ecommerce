package hub

import (
	"bufio"
	"fmt"
	"log"
	"time"
)

type Role struct {
	RoleName string                      `json:"roleName"`
	Clients  map[int64]map[string]Client `json:"clients"`
}

type Client struct {
	Role    string       `json:"role"`
	UUID    string       `json:"uuID"`
	UserID  int64        `json:"userID"`
	Message chan Message `json:"message"`
}

type Message struct {
	Topic   string `json:"topic"`
	Role    string `json:"role"`
	SendID  int64  `json:"sendID"`
	UserID  int64  `json:"userID"`
	Content string `json:"content"`
}

func SendMessage(client Client, w *bufio.Writer, notify <-chan struct{}) {
	ticker := time.NewTicker(5 * time.Second)
	fmt.Fprintf(w, "event: %s\nid: %s:%s\ndata: %s\n\n", "init", client.Role, client.UUID, "hello initial connection")
	w.Flush()

	defer func() {
		log.Println(fmt.Sprintf("SSE connection closed by client %s : %s", client.Role, client.UUID))
		hub.UnRegister <- client
		ticker.Stop()
	}()

	go func() {
		<-notify
	}()

	for loop := true; loop; {
		select {
		case message := <-client.Message:
			fmt.Fprintf(w, "event: %s\nid: %s:%s\ndata: %s\n\n", message.Topic, client.Role, client.UUID, message.Content)
			if err := w.Flush(); err != nil {
				loop = false
			}
		case <-ticker.C:
			fmt.Fprintf(w, "event: %s\nid: %s:%s\ndata: %s\n\n", "ping", client.Role, client.UUID, "pong")
			if err := w.Flush(); err != nil {
				loop = false
			}
		}
	}

}
