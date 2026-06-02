package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// QuickNote holds the schema definition for the QuickNote entity.
type QuickNote struct {
	ent.Schema
}

// Fields of the QuickNote.
func (QuickNote) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.String("content").
			NotEmpty(),
		field.UUID("card_id", uuid.UUID{}),
		field.UUID("created_by", uuid.UUID{}).
			Optional(),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the QuickNote.
func (QuickNote) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("card", Card.Type).
			Ref("quick_notes").
			Unique().
			Required().
			Field("card_id"),
		edge.From("creator", User.Type).
			Ref("quick_notes").
			Unique().
			Field("created_by"),
	}
}
