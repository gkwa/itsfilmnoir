# itsfilmnoir

This project uses Dagger to interact with AWS services, specifically to retrieve the caller identity using the AWS CLI.

## Project Structure

- `dagger/main.go`: Contains the main Dagger module implementation.
- `dagger.json`: Dagger configuration file.
- `justfile`: Contains recipes for common tasks.

## Prerequisites

- [Dagger](https://dagger.io/)
- [Go](https://golang.org/)
- [just](https://github.com/casey/just) command runner
- AWS CLI credentials configured

## Setup

1. Ensure you have Dagger, Go, and just installed on your system.
2. Configure your AWS credentials in `~/.aws/credentials`.

**Important**: The AWS credentials file (`~/.aws/credentials`) is required for this project to function correctly. Make sure it is properly configured with your AWS access key ID and secret access key before running any commands.

## Usage

This project uses `just` as a command runner. Available commands:

```
just                 # List available recipes
just format          # Format Go code and justfile
just get-caller-identity # Execute the GetCallerIdentity function
```

### Get Caller Identity

To retrieve the AWS caller identity:

```
just get-caller-identity
```

This command uses Dagger to create an AWS CLI container, mount your AWS credentials, and execute the `aws sts get-caller-identity` command.

Remember: The AWS credentials file must be present and correctly configured for this command to work.

## Dagger Module

The `Itsfilmnoir` struct in `dagger/main.go` defines the following methods:

- `CreateAWSContainer`: Creates a container with the AWS CLI and mounts AWS credentials.
- `GetCallerIdentity`: Executes the `aws sts get-caller-identity` command in the AWS container.
- `ExecuteGetCallerIdentity`: Combines the above methods to retrieve the caller identity.

## Development

To format the Go code and justfile:

```
just format
```

This command formats the Go code using `gofumpt` and also formats the justfile for consistent styling.

## Configuration

The `dagger.json` file specifies:

- Project name: itsfilmnoir
- SDK: Go

Ensure your Dagger version is compatible with the project configuration.
