package websocket

import (
	"encoding/json"
	"sync"
)

// Hub maintains the set of active clients and broadcasts messages.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// Map of board channels: board_id -> set of clients
	boards map[string]map[*Client]bool

	mu sync.Mutex
}

// NewHub creates a new Hub instance.
func NewHub() *Hub {
	return &Hub{
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		boards:     make(map[string]map[*Client]bool),
	}
}

// Run executes the main Hub event loop.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				// Clean up board memberships
				for boardID, clients := range h.boards {
					if _, exists := clients[client]; exists {
						delete(clients, client)
						h.broadcastToBoardUnsafe(boardID, EventMessage{
							Event: "user_absent",
							Payload: map[string]interface{}{
								"user_id": client.userID.String(),
							},
						})
					}
					if len(clients) == 0 {
						delete(h.boards, boardID)
					}
				}
			}
			h.mu.Unlock()
		}
	}
}

// RegisterClient registers a client.
func (h *Hub) RegisterClient(c *Client) {
	h.register <- c
}

// UnregisterClient unregisters a client.
func (h *Hub) UnregisterClient(c *Client) {
	h.unregister <- c
}

// JoinBoard joins a client to a board channel.
func (h *Hub) JoinBoard(boardID string, client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.boards[boardID] == nil {
		h.boards[boardID] = make(map[*Client]bool)
	}
	h.boards[boardID][client] = true

	// Notify others on the board
	h.broadcastToBoardExceptUnsafe(boardID, EventMessage{
		Event: "user_present",
		Payload: map[string]interface{}{
			"user_id":    client.userID.String(),
			"username":   client.username,
			"avatar_url": client.avatarURL,
		},
	}, client)
}

// LeaveBoard removes a client from a board channel.
func (h *Hub) LeaveBoard(boardID string, client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if clients, ok := h.boards[boardID]; ok {
		if _, exists := clients[client]; exists {
			delete(clients, client)
			h.broadcastToBoardUnsafe(boardID, EventMessage{
				Event: "user_absent",
				Payload: map[string]interface{}{
					"user_id": client.userID.String(),
				},
			})
		}
		if len(clients) == 0 {
			delete(h.boards, boardID)
		}
	}
}

// BroadcastToBoard sends a message to all clients in a board channel.
func (h *Hub) BroadcastToBoard(boardID string, msg EventMessage) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.broadcastToBoardUnsafe(boardID, msg)
}

// BroadcastToBoardExcept sends a message to all clients in a board channel except one.
func (h *Hub) BroadcastToBoardExcept(boardID string, msg EventMessage, except *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.broadcastToBoardExceptUnsafe(boardID, msg, except)
}

func (h *Hub) broadcastToBoardUnsafe(boardID string, msg EventMessage) {
	payload, err := json.Marshal(msg)
	if err != nil {
		return
	}
	if clients, ok := h.boards[boardID]; ok {
		for client := range clients {
			select {
			case client.send <- payload:
			default:
				close(client.send)
				delete(h.clients, client)
				delete(clients, client)
			}
		}
	}
}

func (h *Hub) broadcastToBoardExceptUnsafe(boardID string, msg EventMessage, except *Client) {
	payload, err := json.Marshal(msg)
	if err != nil {
		return
	}
	if clients, ok := h.boards[boardID]; ok {
		for client := range clients {
			if client == except {
				continue
			}
			select {
			case client.send <- payload:
			default:
				close(client.send)
				delete(h.clients, client)
				delete(clients, client)
			}
		}
	}
}

// EventMessage is the standard WebSocket envelope.
type EventMessage struct {
	Event   string                 `json:"event"`
	Payload map[string]interface{} `json:"payload"`
}
