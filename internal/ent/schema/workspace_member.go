package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// WorkspaceMember holds the schema for the workspace_members join table.
type WorkspaceMember struct {
	ent.Schema
}

// Annotations of the WorkspaceMember.
func (WorkspaceMember) Annotations() []schema.Annotation {
	return nil
}

// Fields of the WorkspaceMember.
func (WorkspaceMember) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("workspace_id", uuid.UUID{}),
		field.UUID("user_id", uuid.UUID{}),
		field.Enum("role").
			Values("admin", "member", "viewer").
			Default("member"),
		field.Time("joined_at").
			Default(time.Now).
			Immutable(),
	}
}

// Edges of the WorkspaceMember.
func (WorkspaceMember) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("workspace", Workspace.Type).
			Required().
			Unique().
			Field("workspace_id"),
		edge.To("user", User.Type).
			Required().
			Unique().
			Field("user_id"),
	}
}
