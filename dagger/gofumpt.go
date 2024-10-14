package main

import (
	"context"
	"fmt"

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

	// Execute ls -la before chmod
	lsOutput1, err := containerWithSource.
		WithWorkdir("/src").
		WithExec([]string{"ls", "-la", "dagger/gofumpt.go"}).
		Stdout(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to execute ls -la before chmod: %w", err)
	}
	fmt.Println("ls -la output before chmod:")
	fmt.Println(lsOutput1)

	// Execute chmod
	_, err = containerWithSource.
		WithWorkdir("/src").
		WithExec([]string{"chmod", "-R", "+rw", "."}).
		Sync(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to execute chmod: %w", err)
	}

	// Execute ls -la after chmod
	lsOutput2, err := containerWithSource.
		WithWorkdir("/src").
		WithExec([]string{"ls", "-la", "dagger/gofumpt.go"}).
		Stdout(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to execute ls -la after chmod: %w", err)
	}
	fmt.Println("ls -la output after chmod:")
	fmt.Println(lsOutput2)

	// Execute gofumpt
	output := containerWithSource.
		WithWorkdir("/src").
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
