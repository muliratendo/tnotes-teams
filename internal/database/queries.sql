-- name: GetWorkspaceMemberRole :one
SELECT role FROM workspace_members
WHERE workspace_id = $1 AND user_id = $2;

-- name: ListBoardsByWorkspace :many
SELECT * FROM boards
WHERE workspace_id = $1 AND is_archived = FALSE
ORDER BY created_at DESC;

-- name: GetActivityLogsByBoard :many
SELECT * FROM activity_log
WHERE board_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ShiftCardsPositionAfterInsert :exec
UPDATE cards
SET position = position + 1
WHERE list_id = $1 AND position >= $2;

-- name: ShiftCardsPositionAfterDelete :exec
UPDATE cards
SET position = position - 1
WHERE list_id = $1 AND position > $2;
