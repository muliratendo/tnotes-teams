package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Workspace holds the schema definition for the Workspace entity.
type Workspace struct {
	ent.Schema
}

// Fields of the Workspace.
func (Workspace) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.String("name").
			NotEmpty().
			MaxLen(100),
		field.String("description").
			Optional(),
		field.String("theme_color").
			Default("#6366F1").
			MaxLen(7),
		field.UUID("created_by", uuid.UUID{}).
			Optional(),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
	}
}

// Edges of the Workspace.
func (Workspace) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("boards", Board.Type),
		edge.To("members", User.Type).
			Through("workspace_memberships", WorkspaceMember.Type),
		edge.From("creator", User.Type).
			Ref("created_workspaces").
			Unique().
			Field("created_by"),
	}
}
