package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.String("email").
			Unique().
			NotEmpty().
			MaxLen(255),
		field.String("username").
			Unique().
			NotEmpty().
			MaxLen(50),
		field.String("password_hash").
			NotEmpty().
			Sensitive(),
		field.String("avatar_url").
			Optional().
			MaxLen(500),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("created_workspaces", Workspace.Type),
		edge.To("created_boards", Board.Type),
		edge.To("created_cards", Card.Type),
		edge.To("comments", Comment.Type),
		edge.To("quick_notes", QuickNote.Type),
		edge.To("activity_logs", ActivityLog.Type),
		edge.From("workspaces", Workspace.Type).
			Ref("members").
			Through("workspace_memberships", WorkspaceMember.Type),
	}
}
