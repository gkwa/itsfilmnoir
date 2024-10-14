set shell := ["bash", "-uc"]

default:
    @just --list

get-caller-identity:
    dagger call execute-get-caller-identity --aws-config=file:~/.aws/credentials

format:
    just --unstable --fmt
    dagger call prettier --source=. export --path=.
    dagger call gofumpt --source=. export --path=.

update-versions:
    dagger call update-versions --source=.

update-and-test: update-versions format get-caller-identity
    echo "Update and initial tests complete. Please review changes and run additional tests as needed."
