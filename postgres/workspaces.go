package postgres

import (
	"context"
	"errors"
	"log/slog"

	"github.com/freekobie/hazel/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WorkspaceStore struct {
	conn *pgxpool.Pool
}

func NewWorkspaceStore(conn *pgxpool.Pool) models.WorkspaceStore {
	return &WorkspaceStore{conn: conn}
}

// Create implements models.WorkspaceStore.
func (w *WorkspaceStore) Create(ctx context.Context, ws *models.Workspace) error {
	wsquery := `INSERT INTO workspaces(id, name, description, user_id, created_at, last_modified)
	VALUES($1, $2, $3, $4, $5, $5);`

	memberquery := `INSERT INTO workspace_memberships(workspace_id, user_id)
	VALUES($1, $2)`

	tx, err := w.conn.Begin(ctx)
	if err != nil {
		slog.Error("failed to start transaction", "error", err)
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, wsquery, ws.Id, ws.Name, ws.Description, ws.User.Id, ws.CreatedAt)
	if err != nil {
		slog.Error("failed to create workspace", "error", err)
		return err
	}

	_, err = tx.Exec(ctx, memberquery, ws.Id, ws.User.Id)
	if err != nil {
		slog.Error("failed to create workspace membership", "error", err)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		slog.Error("failed to complete transactions", "error", err)
		return err
	}

	return nil
}

// Get implements models.WorkspaceStore.
func (w *WorkspaceStore) Get(ctx context.Context, id uuid.UUID) (models.Workspace, error) {
	query := `SELECT
	w.id, w.name, w.description, w.created_at, w.last_modified,
	u.id, u.name, u.email, u.profile_picture, u.created_at, u.last_modified, u.verified
	FROM workspaces AS w
	INNER JOIN users AS u
	ON w.user_id = u.id
	WHERE id = $1;`

	row := w.conn.QueryRow(ctx, query, id)

	ws := models.Workspace{}
	err := row.Scan(&ws.Id, &ws.Name, &ws.Description, &ws.CreatedAt, &ws.LastModified, &ws.User.Id, &ws.User.Name, &ws.User.Email, &ws.User.ProfilePhoto, &ws.User.CreatedAt, &ws.User.LastModifed, &ws.User.Verified)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Workspace{}, models.ErrNotFound
		}
		slog.Error("failed to read workspace", "error", err)
		return models.Workspace{}, err
	}

	return ws, nil
}

// GetAllForUser implements models.WorkspaceStore.
func (w *WorkspaceStore) GetAllForUser(ctx context.Context, userId uuid.UUID) ([]models.Workspace, error) {
	query := `SELECT id, name, description, user_id, created_at, last_modified
	FROM workspaces WHERE user_id = $1;`

	rows, err := w.conn.Query(ctx, query, userId)
	if err != nil {
		slog.Error("failed to query rows", "error", err.Error())
		return nil, err
	}

	workspaces := []models.Workspace{}
	for rows.Next() {
		var ws models.Workspace

		err = rows.Scan(&ws.Id, &ws.Name, &ws.Description, &ws.User.Id, &ws.CreatedAt, &ws.LastModified)
		if err != nil {
			slog.Error("failed to scan workspace", "error", err.Error())
			return nil, err
		}

		workspaces = append(workspaces, ws)
	}
	return workspaces, nil
}

// Update implements models.WorkspaceStore.
func (w *WorkspaceStore) Update(ctx context.Context, workspace *models.Workspace) error {
	query := `UPDATE workspaces SET name = $1, description = $2, last_modified = now()
	WHERE id = $3;`

	_, err := w.conn.Exec(ctx, query, workspace.Name, workspace.Description, workspace.Id)
	if err != nil {
		slog.Error("failed to update workspace", "error", err.Error())
		return err
	}

	return nil
}

// Delete implements models.WorkspaceStore.
func (w *WorkspaceStore) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM workspaces WHERE id = $1;`

	_, err := w.conn.Exec(ctx, query, id)

	if err != nil {
		slog.Error("failed to create workspace", "error", err)
		return err
	}

	return nil
}
