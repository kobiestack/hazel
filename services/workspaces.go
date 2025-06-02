package services

import (
	"context"
	"time"

	"github.com/freekobie/hazel/models"
)

type WorkspaceService struct {
	store models.WorkspaceStore
}

func NewWorkspaceService(store models.WorkspaceStore) *WorkspaceService {
	return &WorkspaceService{
		store: store,
	}
}

func (s *WorkspaceService) NewWorkspace(ctx context.Context, ws *models.Workspace) error {
	err := s.store.Create(ctx, ws)
	if err != nil {
		return err
	}

	return nil
}

func (s *WorkspaceService) UpdateWorkspace(ctx context.Context, ws *models.Workspace) error {
	workspace, err := s.store.Get(ctx, ws.Id)
	if err != nil {
		return err
	}

	workspace.Name = ws.Name
	workspace.Description = ws.Description
	workspace.LastModified = time.Now()

	err = s.store.Update(ctx, &workspace)
	if err != nil {
		return err
	}

	return nil
}
