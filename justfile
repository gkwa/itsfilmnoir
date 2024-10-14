set shell := ["bash", "-uc"]

default:
    @just --list

get-caller-identity:
    dagger call execute-get-caller-identity --aws-config=file:~/.aws/credentials

format:
    just --unstable --fmt
    dagger call prettier --source=. export --path=.
