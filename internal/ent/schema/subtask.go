package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Subtask holds the schema definition for the Subtask entity.
type Subtask struct {
	ent.Schema
}

// Fields of the Subtask.
func (Subtask) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.String("title").
			NotEmpty().
			MaxLen(255),
		field.Bool("is_completed").
			Default(false),
		field.Int("position").
			NonNegative(),
		field.UUID("card_id", uuid.UUID{}),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
	}
}

// Edges of the Subtask.
func (Subtask) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("card", Card.Type).
			Ref("subtasks").
			Unique().
			Required().
			Field("card_id"),
	}
}
