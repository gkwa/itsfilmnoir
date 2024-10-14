package main

import (
	"context"

	"dagger/itsfilmnoir/internal/dagger"
)

func (m *Itsfilmnoir) Gofumpt7(ctx context.Context, source *dagger.Directory) (*dagger.Directory, error) {
	gofumptContainer := m.CreateGofumptContainer()

	// Copy the source directory into the container
	containerWithSource := gofumptContainer.WithDirectory("/src", source)

	output := containerWithSource.
		WithWorkdir("/src").
		WithExec([]string{"sudo", "chmod", "-R", "a+rwx", "."}).
		WithExec([]string{"gofumpt", "-w", "."}).
		Directory("/src")

	return output, nil
}

func (m *Itsfilmnoir) Gofumpt7Debug(ctx context.Context, source *dagger.Directory) *dagger.Container {
	gofumptContainer := m.CreateGofumptContainer()
	return gofumptContainer.
		WithMountedDirectory("/src", source).
		WithWorkdir("/src").
		Terminal()
}
