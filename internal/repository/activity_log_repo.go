package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/tendo-mulira/tnotes-teams/internal/ent"
	"github.com/tendo-mulira/tnotes-teams/internal/ent/activitylog"
)

// ActivityLogRepository handles activity log persistence.
type ActivityLogRepository struct {
	client *ent.Client
	db     *sql.DB
}

// NewActivityLogRepository creates a new ActivityLogRepository.
func NewActivityLogRepository(client *ent.Client, db *sql.DB) *ActivityLogRepository {
	return &ActivityLogRepository{client: client, db: db}
}

// Log creates a new activity log entry.
func (r *ActivityLogRepository) Log(ctx context.Context, boardID, userID uuid.UUID, action, entityType string, entityID *uuid.UUID, oldValue, newValue map[string]interface{}) error {
	builder := r.client.ActivityLog.Create().
		SetBoardID(boardID).
		SetUserID(userID).
		SetAction(action).
		SetEntityType(entityType)

	if entityID != nil {
		builder.SetEntityID(*entityID)
	}
	if oldValue != nil {
		builder.SetOldValue(oldValue)
	}
	if newValue != nil {
		builder.SetNewValue(newValue)
	}

	_, err := builder.Save(ctx)
	return err
}

// ListByBoard returns activity logs for a board, newest first.
func (r *ActivityLogRepository) ListByBoard(ctx context.Context, boardID uuid.UUID, limit int) ([]*ent.ActivityLog, error) {
	q := r.client.ActivityLog.Query().
		Where(activitylog.BoardIDEQ(boardID)).
		WithUser().
		Order(ent.Desc(activitylog.FieldCreatedAt))

	if limit > 0 {
		q.Limit(limit)
	}

	return q.All(ctx)
}

// ListByEntity returns activity logs for a specific entity.
func (r *ActivityLogRepository) ListByEntity(ctx context.Context, entityType string, entityID uuid.UUID) ([]*ent.ActivityLog, error) {
	return r.client.ActivityLog.Query().
		Where(
			activitylog.EntityTypeEQ(entityType),
			activitylog.EntityIDEQ(entityID),
		).
		WithUser().
		Order(ent.Desc(activitylog.FieldCreatedAt)).
		All(ctx)
}
