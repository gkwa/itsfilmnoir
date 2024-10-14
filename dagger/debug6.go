package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"dagger/itsfilmnoir/internal/dagger"
)

func (m *Itsfilmnoir) Gofumpt6(ctx context.Context, source *dagger.Directory) (*dagger.Directory, error) {
	gofumptContainer := m.CreateGofumptContainer()
	containerWithSource := gofumptContainer.WithDirectory("/src", source)

	// Create a temporary directory within the container
	tempDir := "/home/linuxbrew/tmp/gofumpt_workspace"
	mkdirCmd := containerWithSource.WithExec([]string{"mkdir", "-p", tempDir})
	_, err := mkdirCmd.Stdout(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	fmt.Printf("Temporary directory created: %s\n", tempDir)

	// Copy source files to the temporary directory
	copyCmd := containerWithSource.WithExec([]string{"sh", "-c", fmt.Sprintf("cp -R /src/. %s && echo 'Copy successful' || echo 'Copy failed'", tempDir)})
	copyOutput, err := copyCmd.Stdout(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to copy source to temp directory: %w", err)
	}
	fmt.Printf("Copy operation result: %s\n", copyOutput)

	// List contents of temp directory
	listTempCmd := containerWithSource.WithExec([]string{"ls", "-la", tempDir})

	// Open an interactive terminal for debugging
	terminalContainer := listTempCmd.Terminal()
	fmt.Println("Opening an interactive terminal for debugging. Type 'exit' when done.")
	_, err = terminalContainer.Stderr(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to open interactive terminal: %w", err)
	}

	// The rest of the function remains the same...
	listTempOutput, err := listTempCmd.Stdout(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list temp directory contents: %w", err)
	}
	fmt.Printf("Temp directory contents:\n%s\n", listTempOutput)

	// Count number of files in the temporary directory
	countCmd := containerWithSource.WithExec([]string{"sh", "-c", fmt.Sprintf("find %s -type f | wc -l", tempDir)})
	countOutput, err := countCmd.Stdout(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count files in temp directory: %w", err)
	}
	fileCount, err := strconv.Atoi(strings.TrimSpace(countOutput))
	if err != nil {
		return nil, fmt.Errorf("failed to parse file count: %w", err)
	}
	fmt.Printf("Number of files in temp directory: %d\n", fileCount)
	if fileCount < 2 {
		return nil, fmt.Errorf("insufficient files in temp directory: expected at least 2, got %d", fileCount)
	}

	// Run gofumpt
	gofumptCmd := containerWithSource.
		WithWorkdir(tempDir).
		WithExec([]string{"gofumpt", "-w", "."})

	gofumptOutput, err := gofumptCmd.Stdout(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to run gofumpt: %w", err)
	}
	fmt.Println("gofumpt output:")
	fmt.Println(gofumptOutput)

	return containerWithSource.Directory(tempDir), nil
}
