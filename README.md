# GoLuca

A Simple Accounting Ledger


[![Go Report Card](https://goreportcard.com/badge/github.com/hampgoodwin/GoLuca)](https://goreportcard.com/report/github.com/hampgoodwin/GoLuca) [![Coverage Status](https://coveralls.io/repos/github/hampgoodwin/GoLuca/badge.svg)](https://coveralls.io/github/hampgoodwin/GoLuca) [![golangci-lint](https://github.com/hampgoodwin/GoLuca/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/hampgoodwin/GoLuca/actions/workflows/golangci-lint.yml)
[![buf-lint](https://github.com/hampgoodwin/GoLuca/actions/workflows/buf-lint.yml/badge.svg)](https://github.com/hampgoodwin/GoLuca/actions/workflows/buf-lint.yml)

- Simple application which writes and reads accounting ledger entries

## Tooling

Core dependencies are nix and nix-direnv. Everything else required for this repository is managed by these two!

We use [direnv](https://direnv.net/) to manage all toolings and dependencies. You can reference the ./flake.nix file for their declarations and ./flake.lock for versions. Currently we are not pinning versions in flake.nix via import, but we should do that soon. Below are some tools we use, and a description of how we use them.

- [nix-direnv](https://github.com/nix-community/nix-direnv)
    - we use direnv, as started above, but we use a specific nix-direnv that allows us to use nix, specifically nix develop via flake.nix, to manage our toolings such that we're all using the same thing.
- [buf](https://buf.build/)
    - we use buf for proto toolchain. Linting, formatting, and generation. Also in CI allows for backwards-breaking change detection.
- [gofumpt](https://github.com/mvdan/gofumpt) for formatting.
    - gofumpy for formatting go code.
- [golangci-lint](https://github.com/golangci/golangci-lint) for linting.
    - golangci-lint for linting go code. I like to update golangci-lint versions and use it's default opinions; hence we have no .golangci-lint configuration file.
- [colima](https://github.com/abiosoft/colima) for container runtimes.
- [jaeger](https://www.jaegertracing.io/) for local trace collector and ui.
- [nats](https://nats.io/) for eventing.

## How to develop

Read the [Tooling](##tooling) section.

Navigate in your shell to this repos directory and a developer environment will be configured for you. Either have a global container runtime available, or use this environments colima, by running `colima start`. Review the [Makefile](./Makefile) to see available commands. Otherwise, everything is what you'd expect. Best of luck!

