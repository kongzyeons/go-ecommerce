package hub

import (
	"context"
	"fmt"
	"log"
	"sync"
)

// hub : {
// "admin" : {
// 	userID : {
// 		uuid : struct
// 	}
// }
// "user" : {
// 	userID : {
// 		uuid : struct
// 	}}
// }

type Hub struct {
	Roles      map[string]Role
	Register   chan Client
	UnRegister chan Client
	Broadcast  chan Message
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
	mu         sync.RWMutex
}

var hub *Hub
var hubOnce sync.Once

func InitHub() *Hub {
	hubOnce.Do(func() {
		ctx, cancel := context.WithCancel(context.Background())
		hub = &Hub{
			Roles:      make(map[string]Role),
			Register:   make(chan Client),
			UnRegister: make(chan Client),
			Broadcast:  make(chan Message, 100),
			ctx:        ctx,
			cancel:     cancel,
		}
	})
	return hub
}

func CloseHub() {
	if hub != nil {
		hub.Stop()
	}
}

func (h *Hub) Stop() {
	h.cancel()
	h.wg.Wait()
	h.mu.Lock()
	defer h.mu.Unlock()
	for c := range h.Roles {
		for _, client := range h.Roles[c].Clients {
			for _, id := range client {
				close(id.Message)
			}
		}
		delete(h.Roles, c)
	}
}

func (h *Hub) Run() {
	log.Println("Hub started")
	h.wg.Add(1)
	go func() {
		defer h.wg.Done()
		for {
			select {
			case <-h.ctx.Done():
				log.Println("Hub stopped")
				return
			case c := <-h.Register:
				h.mu.Lock()
				if _, ok := h.Roles[c.Role]; !ok {
					h.Roles[c.Role] = Role{
						RoleName: c.Role,

						Clients: make(map[int64]map[string]Client),
					}
				}
				if _, ok := h.Roles[c.Role].Clients[c.UserID]; !ok {
					h.Roles[c.Role].Clients[c.UserID] = make(map[string]Client)
				}
				h.Roles[c.Role].Clients[c.UserID][c.UUID] = c
				log.Println(fmt.Sprintf("register lenght role : %d", len(h.Roles)))
				log.Println(fmt.Sprintf("register lenght role %s : %d", c.Role, len(h.Roles[c.Role].Clients)))
				h.mu.Unlock()

			case c := <-h.UnRegister:
				h.mu.Lock()
				if _, ok := h.Roles[c.Role]; ok {
					if _, ok := h.Roles[c.Role].Clients[c.UserID]; ok {
						delete(h.Roles[c.Role].Clients[c.UserID], c.UUID)
						if len(h.Roles[c.Role].Clients[c.UserID]) == 0 {
							delete(h.Roles[c.Role].Clients, c.UserID)
						}
					}
					if len(h.Roles[c.Role].Clients) == 0 {
						delete(h.Roles, c.Role)
					}
				}
				log.Println(fmt.Sprintf("unrister lenght role : %d", len(h.Roles)))
				log.Println(fmt.Sprintf("unrister lenght role %s : %d", c.Role, len(h.Roles[c.Role].Clients)))

			case msg := <-h.Broadcast:
				if _, ok := h.Roles[msg.Role]; ok {
					if _, ok := h.Roles[msg.Role].Clients[msg.SendID]; ok {
						for _, client := range h.Roles[msg.Role].Clients[msg.SendID] {
							client.Message <- msg
						}
					}
				}
			}
		}
	}()
}
