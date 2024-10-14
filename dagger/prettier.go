package main

import (
	"context"

	"dagger/itsfilmnoir/internal/dagger"
)

func (m *Itsfilmnoir) CreatePrettierContainer(ctx context.Context) (*dagger.Container, error) {
	nodeCache := dag.CacheVolume("node")
	container := dag.Container().
		From("node:lts-alpine3.20").
		WithMountedCache("/root/.npm", nodeCache).
		WithExec([]string{"npm", "install", "-g", "prettier"})
	return container, nil
}

func (m *Itsfilmnoir) PrettierDebug(ctx context.Context, source *dagger.Directory) *dagger.Container {
	prettierContainer, err := m.CreatePrettierContainer(ctx)
	if err != nil {
		return dag.Container().WithExec([]string{"echo", "Error creating Prettier container: " + err.Error()})
	}
	return prettierContainer.
		WithMountedDirectory("/src", source).
		WithWorkdir("/src").
		Terminal()
}

func (m *Itsfilmnoir) Prettier(ctx context.Context, source *dagger.Directory) (*dagger.Directory, error) {
	prettierContainer, err := m.CreatePrettierContainer(ctx)
	if err != nil {
		return nil, err
	}
	output := prettierContainer.
		WithMountedDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"prettier", "--write", "."}).
		Directory("/src")
	return output, nil
}
