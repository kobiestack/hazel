package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/freekobie/hazel/models"
	"github.com/freekobie/hazel/postgres"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestWorkspace(name string, user *models.User) *models.Workspace {
	return &models.Workspace{
		Id:           uuid.New(),
		Name:         name,
		Description:  "Test workspace description",
		User:         user,
		CreatedAt:    time.Now().UTC(),
		LastModified: time.Now().UTC(),
	}
}

func TestWorkspaceStore_Create(t *testing.T) {
	pool := setupTestDB(t)
	workspaceStore := postgres.NewWorkspaceStore(pool)
	userStore := postgres.NewUserStore(pool)
	ctx := context.Background()

	user := createTestUser("Test User", generateTestEmail())
	require.NoError(t, userStore.InsertUser(ctx, user))

	tests := []struct {
		name      string
		workspace *models.Workspace
		wantErr   bool
	}{
		{
			name:      "valid workspace",
			workspace: createTestWorkspace("Test Workspace", user),
			wantErr:   false,
		},
		{
			name: "empty name",
			workspace: &models.Workspace{
				Id:          uuid.New(),
				Description: "Test description",
				User:        user,
				CreatedAt:   time.Now().UTC(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := workspaceStore.Create(ctx, tt.workspace)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			got, err := workspaceStore.Get(ctx, tt.workspace.Id)
			require.NoError(t, err)
			assert.Equal(t, tt.workspace.Id, got.Id)
			assert.Equal(t, tt.workspace.Name, got.Name)
			assert.Equal(t, tt.workspace.Description, got.Description)
		})
	}
}

func TestWorkspaceStore_Get(t *testing.T) {
	pool := setupTestDB(t)
	workspaceStore := postgres.NewWorkspaceStore(pool)
	userStore := postgres.NewUserStore(pool)
	ctx := context.Background()

	user := createTestUser("Test User", generateTestEmail())
	require.NoError(t, userStore.InsertUser(ctx, user))

	workspace := createTestWorkspace("Get Test", user)
	require.NoError(t, workspaceStore.Create(ctx, workspace))

	tests := []struct {
		name    string
		id      uuid.UUID
		wantErr bool
	}{
		{
			name:    "existing workspace",
			id:      workspace.Id,
			wantErr: false,
		},
		{
			name:    "non-existent workspace",
			id:      uuid.New(),
			wantErr: true,
		},
		{
			name:    "invalid ID",
			id:      uuid.Nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := workspaceStore.Get(ctx, tt.id)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, workspace.Id, got.Id)
			assert.Equal(t, workspace.Name, got.Name)
			assert.Equal(t, workspace.Description, got.Description)
		})
	}
}

func TestWorkspaceStore_Delete(t *testing.T) {
	pool := setupTestDB(t)
	workspaceStore := postgres.NewWorkspaceStore(pool)
	userStore := postgres.NewUserStore(pool)
	ctx := context.Background()

	user := createTestUser("Test User", generateTestEmail())
	require.NoError(t, userStore.InsertUser(ctx, user))

	workspace := createTestWorkspace("Delete Test", user)
	require.NoError(t, workspaceStore.Create(ctx, workspace))

	tests := []struct {
		name    string
		id      uuid.UUID
		wantErr bool
	}{
		{
			name:    "existing workspace",
			id:      workspace.Id,
			wantErr: false,
		},
		{
			name:    "non-existent workspace",
			id:      uuid.New(),
			wantErr: true,
		},
		{
			name:    "invalid ID",
			id:      uuid.Nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := workspaceStore.Delete(ctx, tt.id)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			_, err = workspaceStore.Get(ctx, tt.id)
			assert.Error(t, err)
		})
	}
}
