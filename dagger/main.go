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

func (m *Itsfilmnoir) Prettier(ctx context.Context, source *dagger.Directory) error {
	prettierContainer, err := m.CreatePrettierContainer(ctx)
	if err != nil {
		return err
	}

	_, err = prettierContainer.
		WithMountedDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"prettier", "--write", "."}).
		Sync(ctx)

	return err
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
