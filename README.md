# itsfilmnoir

Practice using Dagger with secrets in AWS environment.

## Motivation

This project serves as a practical exercise for deploying secrets to a container using Dagger and verifying that the AWS CLI can detect and use these secrets.

## Goal

To demonstrate how to securely handle AWS credentials with Dagger, allowing you to run AWS CLI commands in a containerized environment.

## Prerequisites

- Dagger
- just command runner
- AWS CLI credentials

## Quick Start

1. Ensure your AWS credentials are in `~/.aws/credentials` on your docker host.

2. Run the caller identity check:
   ```
   just get-caller-identity
   ```

This command creates an AWS CLI container, mounts your credentials, and runs `aws sts get-caller-identity`.

## Available Commands

```
just                 # List available recipes
just format          # Format Go code and justfile
just get-caller-identity # Run AWS caller identity check
```

## Project Structure

- `dagger/main.go`: Dagger module implementation
- `dagger.json`: Dagger configuration
- `justfile`: Command recipes

For more details on the implementation, check the comments in `dagger/main.go`.
