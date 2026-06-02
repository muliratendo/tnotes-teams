package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// List holds the schema definition for the List (column) entity.
type List struct {
	ent.Schema
}

// Fields of the List.
func (List) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.String("title").
			NotEmpty().
			MaxLen(100),
		field.Int("position").
			NonNegative(),
		field.UUID("board_id", uuid.UUID{}),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
	}
}

// Edges of the List.
func (List) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("cards", Card.Type),
		edge.From("board", Board.Type).
			Ref("lists").
			Unique().
			Required().
			Field("board_id"),
	}
}
