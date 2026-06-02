package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// ActivityLog holds the schema definition for the ActivityLog entity.
type ActivityLog struct {
	ent.Schema
}

// Fields of the ActivityLog.
func (ActivityLog) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.String("action").
			NotEmpty().
			MaxLen(50),
		field.String("entity_type").
			Optional().
			MaxLen(50),
		field.UUID("entity_id", uuid.UUID{}).
			Optional().
			Nillable(),
		field.JSON("old_value", map[string]interface{}{}).
			Optional(),
		field.JSON("new_value", map[string]interface{}{}).
			Optional(),
		field.UUID("board_id", uuid.UUID{}),
		field.UUID("user_id", uuid.UUID{}).
			Optional(),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
	}
}

// Edges of the ActivityLog.
func (ActivityLog) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("board", Board.Type).
			Ref("activity_logs").
			Unique().
			Required().
			Field("board_id"),
		edge.From("user", User.Type).
			Ref("activity_logs").
			Unique().
			Field("user_id"),
	}
}
