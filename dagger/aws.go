package main

import (
	"context"

	"dagger/itsfilmnoir/internal/dagger"
)

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
