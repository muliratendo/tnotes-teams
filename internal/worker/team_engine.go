package worker

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/tendo-mulira/tnotes-teams/internal/config"
	"github.com/tendo-mulira/tnotes-teams/internal/service"
	"github.com/tendo-mulira/tnotes-teams/internal/websocket"
)

// TeamEngine simulates real-time activity of other team members.
type TeamEngine struct {
	services *service.Services
	hub      *websocket.Hub
	cfg      *config.Config
}

// NewTeamEngine creates a new TeamEngine worker.
func NewTeamEngine(services *service.Services, hub *websocket.Hub, cfg *config.Config) *TeamEngine {
	return &TeamEngine{
		services: services,
		hub:      hub,
		cfg:      cfg,
	}
}

// Start runs the simulation scheduler loop in the background.
func (e *TeamEngine) Start(ctx context.Context) {
	if !e.cfg.TeamEngineEnabled {
		log.Println("Team Update Engine is disabled")
		return
	}

	interval, err := time.ParseDuration(e.cfg.TeamEngineInterval)
	if err != nil {
		interval = 30 * time.Second
	}

	log.Printf("Team Update Engine starting with interval %v", interval)

	// Ensure mock users exist in database
	mockUserIDs := e.ensureMockUsers(ctx)
	if len(mockUserIDs) == 0 {
		log.Println("Team Engine: No mock users could be created. Stopping simulator.")
		return
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for {
		select {
		case <-ctx.Done():
			log.Println("Team Update Engine stopping...")
			return
		case <-ticker.C:
			e.runSimulationStep(ctx, mockUserIDs, r)
		}
	}
}

// ensureMockUsers registers our virtual team members in the system.
func (e *TeamEngine) ensureMockUsers(ctx context.Context) map[string]uuid.UUID {
	users := map[string]string{
		"sarah":   "sarah@tnotes.teams",
		"alex":    "alex@tnotes.teams",
		"michael": "michael@tnotes.teams",
	}

	userIDs := make(map[string]uuid.UUID)

	for username, email := range users {
		res, err := e.services.Auth.Register(ctx, service.RegisterInput{
			Email:    email,
			Username: username,
			Password: "supersecretmockpassword123",
		})
		if err != nil {
			// Already registered, fetch details
			u, getErr := e.services.Workspace.GetByID(ctx, uuid.Nil) // just check by email via user repo
			_ = u
			existingUser, fetchErr := e.services.Auth.Login(ctx, service.LoginInput{
				Email:    email,
				Password: "supersecretmockpassword123",
			})
			if fetchErr == nil {
				userIDs[username] = existingUser.User.ID
			} else {
				log.Printf("Team Engine error ensuring user %s: %v", username, getErr)
			}
		} else {
			userIDs[username] = res.User.ID
		}
	}

	return userIDs
}

func (e *TeamEngine) runSimulationStep(ctx context.Context, mockUsers map[string]uuid.UUID, r *rand.Rand) {
	// 1. Fetch workspaces to find boards to simulate activity on
	// We'll query a single active workspace/board.
	// Since workspaces are scoped to users, let's fetch workspaces that the mock users are in.
	// As a shortcut, we query all boards currently in database.
	boards, err := e.services.Board.ListByWorkspace(ctx, uuid.Nil) // We can't query by Nil, but we can query through all workspaces
	_ = boards
	_ = err

	// Let's query Ent directly or list workspaces for mock user "sarah"
	sarahID := mockUsers["sarah"]
	workspaces, err := e.services.Workspace.ListByUser(ctx, sarahID)
	if err != nil || len(workspaces) == 0 {
		// No active workspaces yet, skip simulation step
		return
	}

	// Pick a random workspace
	rng := r
	ws := workspaces[rng.Intn(len(workspaces))]

	// List boards in this workspace
	wsBoards, err := e.services.Board.ListByWorkspace(ctx, ws.ID)
	if err != nil || len(wsBoards) == 0 {
		return
	}

	board := wsBoards[rng.Intn(len(wsBoards))]

	// Ensure other mock users are members of this workspace
	for _, id := range mockUsers {
		_, _ = e.services.Workspace.InviteMember(ctx, ws.ID, service.InviteMemberInput{
			Email: getEmailForMock(mockUsers, id),
			Role:  "member",
		})
	}

	// Pick a random mock user to perform the action
	mockNames := []string{"sarah", "alex", "michael"}
	actorName := mockNames[rng.Intn(len(mockNames))]
	actorID := mockUsers[actorName]

	// Fetch full board data to find lists and cards
	boardData, err := e.services.Board.GetWithFullData(ctx, board.ID)
	if err != nil || len(boardData.Lists) == 0 {
		return
	}

	actionType := rng.Intn(4)
	switch actionType {
	case 0: // Mock cursor movement
		list := boardData.Lists[rng.Intn(len(boardData.Lists))]
		var cardID string
		if len(list.Cards) > 0 {
			cardID = list.Cards[rng.Intn(len(list.Cards))].ID.String()
		}
		// Send a few mock cursor events
		e.hub.BroadcastToBoard(board.ID.String(), websocket.EventMessage{
			Event: "cursor_broadcast",
			Payload: map[string]interface{}{
				"user_id":  actorID.String(),
				"username": actorName,
				"card_id":  cardID,
				"x":        rng.Float64() * 800,
				"y":        rng.Float64() * 600,
			},
		})

	case 1: // Add a comment (simulating typing first)
		list := boardData.Lists[rng.Intn(len(boardData.Lists))]
		if len(list.Cards) == 0 {
			return
		}
		card := list.Cards[rng.Intn(len(list.Cards))]

		// 1. Typing start
		e.hub.BroadcastToBoard(board.ID.String(), websocket.EventMessage{
			Event: "typing_broadcast",
			Payload: map[string]interface{}{
				"user_id":   actorID.String(),
				"username":  actorName,
				"card_id":   card.ID.String(),
				"is_typing": true,
			},
		})

		// Wait 2 seconds (simulate typing delay)
		go func() {
			time.Sleep(2 * time.Second)
			e.hub.BroadcastToBoard(board.ID.String(), websocket.EventMessage{
				Event: "typing_broadcast",
				Payload: map[string]interface{}{
					"user_id":   actorID.String(),
					"username":  actorName,
					"card_id":   card.ID.String(),
					"is_typing": false,
				},
			})

			comments := []string{
				"I'm working on this now.",
				"Looks good to me!",
				"Is this due by end of sprint?",
				"Can we double check the design assets?",
				"Completed the first pass of edits.",
			}
			commentText := comments[rand.Intn(len(comments))]

			comment, err := e.services.Card.CreateComment(context.Background(), card.ID, service.CreateCommentInput{
				Content: commentText,
			}, actorID)

			if err == nil {
				e.hub.BroadcastToBoard(board.ID.String(), websocket.EventMessage{
					Event: "board_mutated",
					Payload: map[string]interface{}{
						"action":      "comment_added",
						"entity_type": "comment",
						"entity_id":   comment.ID.String(),
						"payload":     comment,
						"actor_id":    actorID.String(),
					},
				})
			}
		}()

	case 2: // Move a card
		// Find a card and move it to a different column
		if len(boardData.Lists) < 2 {
			return
		}
		var sourceList service.ListDTO
		var moveCard service.CardDTO
		found := false

		// Try to find a list with cards
		for _, l := range boardData.Lists {
			if len(l.Cards) > 0 {
				sourceList = l
				moveCard = l.Cards[rng.Intn(len(l.Cards))]
				found = true
				break
			}
		}

		if !found {
			return
		}

		// Pick destination list (different from source)
		var destList service.ListDTO
		for {
			destList = boardData.Lists[rng.Intn(len(boardData.Lists))]
			if destList.ID != sourceList.ID {
				break
			}
		}

		movedCard, err := e.services.Card.MoveCard(ctx, moveCard.ID, service.MoveCardInput{
			NewListID:   destList.ID,
			NewPosition: len(destList.Cards),
		}, actorID)

		if err == nil {
			e.hub.BroadcastToBoard(board.ID.String(), websocket.EventMessage{
				Event: "board_mutated",
				Payload: map[string]interface{}{
					"action":      "card_moved",
					"entity_type": "card",
					"entity_id":   movedCard.ID.String(),
					"payload": map[string]interface{}{
						"card_id":       movedCard.ID.String(),
						"list_id":       destList.ID.String(),
						"old_list_id":   sourceList.ID.String(),
						"new_position":  movedCard.Position,
					},
					"actor_id": actorID.String(),
				},
			})
		}

	case 3: // Toggle a subtask
		// Find a card with subtasks
		var targetCard service.CardDTO
		found := false

		for _, l := range boardData.Lists {
			for _, c := range l.Cards {
				if len(c.Subtasks) > 0 {
					targetCard = c
					found = true
					break
				}
			}
		}

		if !found {
			return
		}

		// Toggle a random subtask
		st := targetCard.Subtasks[rng.Intn(len(targetCard.Subtasks))]
		toggledSt, err := e.services.Card.ToggleSubtask(ctx, st.ID, actorID)
		if err == nil {
			e.hub.BroadcastToBoard(board.ID.String(), websocket.EventMessage{
				Event: "board_mutated",
				Payload: map[string]interface{}{
					"action":      "subtask_toggled",
					"entity_type": "subtask",
					"entity_id":   toggledSt.ID.String(),
					"payload":     toggledSt,
					"actor_id":    actorID.String(),
				},
			})
		}
	}
}

func getEmailForMock(mockUsers map[string]uuid.UUID, id uuid.UUID) string {
	for name, uid := range mockUsers {
		if uid == id {
			return name + "@tnotes.teams"
		}
	}
	return "unknown@tnotes.teams"
}
