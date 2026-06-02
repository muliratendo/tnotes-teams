package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Card holds the schema definition for the Card entity.
type Card struct {
	ent.Schema
}

// Fields of the Card.
func (Card) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.String("title").
			NotEmpty().
			MaxLen(255),
		field.String("description").
			Optional(),
		field.Int("position").
			NonNegative(),
		field.Time("due_date").
			Optional().
			Nillable(),
		field.JSON("labels", []string{}).
			Optional(),
		field.Int("progress_percentage").
			Default(0).
			Min(0).
			Max(100),
		field.UUID("list_id", uuid.UUID{}),
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

// Edges of the Card.
func (Card) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("subtasks", Subtask.Type),
		edge.To("quick_notes", QuickNote.Type),
		edge.To("comments", Comment.Type),
		edge.From("list", List.Type).
			Ref("cards").
			Unique().
			Required().
			Field("list_id"),
		edge.From("creator", User.Type).
			Ref("created_cards").
			Unique().
			Field("created_by"),
	}
}
