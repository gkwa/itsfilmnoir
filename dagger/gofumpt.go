package main

import (
	"context"

	"dagger/itsfilmnoir/internal/dagger"
)

func (m *Itsfilmnoir) CreateGofumptContainer() *dagger.Container {
	return dag.Container().
		From("homebrew/brew:4.4.1").
		WithExec([]string{"brew", "install", "gofumpt"})
}

func (m *Itsfilmnoir) Gofumpt(ctx context.Context, source *dagger.Directory) (*dagger.Directory, error) {
	gofumptContainer := m.CreateGofumptContainer()

	// Copy the source directory into the container, excluding .git
	dirOptions := dagger.ContainerWithDirectoryOpts{
		Exclude: []string{".git"},
	}
	containerWithSource := gofumptContainer.WithDirectory("/src", source, dirOptions)

	output := containerWithSource.
		WithWorkdir("/src").
		WithExec([]string{"sudo", "chmod", "-R", "a+rwx", "."}).
		WithExec([]string{"gofumpt", "-w", "."}).
		Directory("/src")

	return output, nil
}

func (m *Itsfilmnoir) GofumptDebug(ctx context.Context, source *dagger.Directory) *dagger.Container {
	gofumptContainer := m.CreateGofumptContainer()
	return gofumptContainer.
		WithMountedDirectory("/src", source).
		WithWorkdir("/src").
		Terminal()
}
