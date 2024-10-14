package main

import (
	"context"

	"dagger/itsfilmnoir/internal/dagger"
)

type Itsfilmnoir struct{}

func (m *Itsfilmnoir) CreatePrettierContainer(ctx context.Context) (*dagger.Container, error) {
	nodeCache := dag.CacheVolume("node")
	container := dag.Container().
		From("node:latest").
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

func (m *Itsfilmnoir) CreateGofumptContainer() *dagger.Container {
	return dag.Container().
		From("homebrew/brew").
		WithExec([]string{"brew", "install", "gofumpt"})
}

func (m *Itsfilmnoir) Gofumpt(ctx context.Context, source *dagger.Directory) (*dagger.Directory, error) {
	gofumptContainer := m.CreateGofumptContainer()
	
	// Copy the source directory into the container
	containerWithSource := gofumptContainer.WithDirectory("/src", source)

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







func (m *Itsfilmnoir) CreateAWSContainer(awsConfig *dagger.Secret) *dagger.Container {
	return dag.Container().
		From("amazon/aws-cli:latest").
		WithMountedSecret("/root/.aws/credentials", awsConfig)
}

func (m *Itsfilmnoir) GetCallerIdentity(ctx context.Context, awsContainer *dagger.Container) (string, error) {
	output, err := awsContainer.
		WithExec([]string{"aws", "sts", "get-caller-identity"}).Stdout(ctx)
	if err != nil {
		return "", err
	}
	return output, nil
}

func (m *Itsfilmnoir) ExecuteGetCallerIdentity(ctx context.Context, awsConfig *dagger.Secret) (string, error) {
	awsContainer := m.CreateAWSContainer(awsConfig)
	return m.GetCallerIdentity(ctx, awsContainer)
}


func (m *Itsfilmnoir) GofumptDebug2(ctx context.Context, source *dagger.Directory) *dagger.Container {
    gofumptContainer := m.CreateGofumptContainer()

    // Copy the source directory to a temp directory inside the container, work there
    return gofumptContainer.
        WithMountedDirectory("/src", source). // Mount the source directory
        WithExec([]string{"cp", "-r", "/src", "/tmp/src"}). // Copy files to /tmp/src
        WithWorkdir("/tmp/src"). // Set working directory to /tmp/src
        WithExec([]string{"touch", "/tmp/src/testfile.txt"}). // Test write access by creating a file in /tmp/src
        Terminal()
}


func (m *Itsfilmnoir) GofumptDebug3(ctx context.Context, source *dagger.Directory) *dagger.Container {
    gofumptContainer := m.CreateGofumptContainer()

    // Mount the directory, copy as root, then switch back to user 1000:1000
    return gofumptContainer.
        WithMountedDirectory("/src", source). // Mount source directory
        WithUser("0:0"). // Temporarily switch to root to copy files
        WithExec([]string{"cp", "-r", "/src", "/tmp/src"}). // Copy files to /tmp/src
        WithUser("1000:1000"). // Switch back to user 1000:1000
        WithWorkdir("/tmp/src"). // Set working directory to /tmp/src
        WithExec([]string{"touch", "/tmp/src/testfile.txt"}). // Test write access in /tmp/src
        Terminal()
}


func (m *Itsfilmnoir) GofumptDebug4(ctx context.Context, source *dagger.Directory) *dagger.Container {
    gofumptContainer := m.CreateGofumptContainer()

    // Mount the directory and preserve file modes while copying
    return gofumptContainer.
        WithMountedDirectory("/src", source). // Mount source directory
        WithExec([]string{"cp", "-rp", "/src", "/tmp/src"}). // Copy files with permissions preserved
        WithWorkdir("/tmp/src"). // Set working directory to /tmp/src
        WithExec([]string{"touch", "/tmp/src/testfile.txt"}). // Test write access in /tmp/src
        Terminal()
}