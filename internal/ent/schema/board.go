package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Board holds the schema definition for the Board entity.
type Board struct {
	ent.Schema
}

// Fields of the Board.
func (Board) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.String("name").
			NotEmpty().
			MaxLen(100),
		field.String("description").
			Optional(),
		field.String("color_theme").
			Default("#6366F1").
			MaxLen(7),
		field.Bool("is_archived").
			Default(false),
		field.UUID("workspace_id", uuid.UUID{}),
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

// Edges of the Board.
func (Board) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("lists", List.Type),
		edge.To("activity_logs", ActivityLog.Type),
		edge.From("workspace", Workspace.Type).
			Ref("boards").
			Unique().
			Required().
			Field("workspace_id"),
		edge.From("creator", User.Type).
			Ref("created_boards").
			Unique().
			Field("created_by"),
	}
}
