package websocket

// routeEvent routes the incoming client event message to appropriate handlers or hub broadcasts.
func (c *Client) routeEvent(env EventMessage) {
	switch env.Event {
	case "join_board":
		boardID, _ := env.Payload["board_id"].(string)
		if boardID != "" {
			c.hub.JoinBoard(boardID, c)
		}
	case "leave_board":
		boardID, _ := env.Payload["board_id"].(string)
		if boardID != "" {
			c.hub.LeaveBoard(boardID, c)
		}
	case "cursor_update":
		boardID, _ := env.Payload["board_id"].(string)
		cardID, _ := env.Payload["card_id"].(string)
		x, _ := env.Payload["x"].(float64)
		y, _ := env.Payload["y"].(float64)
		if boardID != "" {
			c.hub.BroadcastToBoardExcept(boardID, EventMessage{
				Event: "cursor_broadcast",
				Payload: map[string]interface{}{
					"user_id":  c.userID.String(),
					"username": c.username,
					"card_id":  cardID,
					"x":        x,
					"y":        y,
				},
			}, c)
		}
	case "typing_start":
		boardID, _ := env.Payload["board_id"].(string)
		cardID, _ := env.Payload["card_id"].(string)
		if boardID != "" && cardID != "" {
			c.hub.BroadcastToBoardExcept(boardID, EventMessage{
				Event: "typing_broadcast",
				Payload: map[string]interface{}{
					"user_id":   c.userID.String(),
					"username":  c.username,
					"card_id":   cardID,
					"is_typing": true,
				},
			}, c)
		}
	case "typing_stop":
		boardID, _ := env.Payload["board_id"].(string)
		cardID, _ := env.Payload["card_id"].(string)
		if boardID != "" && cardID != "" {
			c.hub.BroadcastToBoardExcept(boardID, EventMessage{
				Event: "typing_broadcast",
				Payload: map[string]interface{}{
					"user_id":   c.userID.String(),
					"username":  c.username,
					"card_id":   cardID,
					"is_typing": false,
				},
			}, c)
		}
	}
}
