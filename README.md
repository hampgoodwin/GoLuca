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

---

TODO
- [x] use nix devlopment environment
  - I was in the middle of updating all of my nvim configurations, lsp, lint, format etc.. in the midst of this I also updated my make targets for running my linter golangci-lint in a container. It then also occurred to me that I run linter in my ide, and I would need to manually manage a global golangci-lint version and ensure it's the same as my ci flow. That sounds like a bit of a pain and also inflexible between different repo's packages. For just linting it's probably not a big deal, but there are other tools where it could be more important/impactful. It occurred to me to use something like a multiple version manager, such as asdf. This would definitely work, but then there is another concern of managing the build tool versions separate from that.
  - I remember hearing about devenv.sh, and briefly about nix develop/shell. I did a _little_ bit of research and it seems like nix development via defining a nix.flake will get me not only consist reproducible distinct projct builds but gains for my ci and build process. SO, I'm pretty interested in this as it is a place I could get lots of gains for developer productivity.
- [ ] update various patterns/paradigms
  - [ ] singleton logger
  - [ ] where interfaces are defined
  - [ ] update to proto (editions)[https://protobuf.dev/editions/overview/#migrating]...(lang guide)[https://protobuf.dev/programming-guides/editions]
  - [ ] update pagination, make a separate hampgoodwin/go-paginate lib
  - [ ] update testing paradigm
  - [ ] remove http interface, grpc-gateway, proto validation?
  - [ ] swap ksuid for uuidv7!!
- [x] update all dependencies
- [x] make a nix dev environment configuration
- [ ] after f/e deploy this thing somewhere
- [ ] create tenant data structures
  - [ ] create auth
  - [ ] create management ui
- [ ] add govulnscan
- [ ] grpc-gateway?
- [ ] o11y
    - [X] Add tracing
      - [x] Add tracing to response.go files
    - [ ] Add metrics
        - [ ] OpenTelemetry metrics ready?
    - [ ] Log Collector?
    - [ ] move off Jaeger, and all in on Grafana stack
- [ ] make a frontend with some dashboard functionality, (vue3 plz) OOOOH!!!
- [ ] replace hampgoodwin/errors with go std lib multi-errors <3!!!
- [ ] health checks
    - [ ] add a health reporter
    - [ ] add health checks
- [ ] improve pagination by displaying page number of result.
- [ ] Testing
    - [x] Add more unit tests, solitary and sociable
    - [x] Add Integration testing
    - [x] Add fail case account controller testing
    - [x] Use httptest for tests instead of setting up an http server for each test
        - [x] Fix bug where connection is rejected (I sus due to db+http connection limitations in unix)
    - [x] Improve test coverage by testing more grpc methods
    - [ ] Add cursor testing for listing accounts
    - [ ] Add fuzzing
    - [ ] Create a seeder for a basic dev environment of data
    - [ ] Do some stress testing; how much data and we throw at & get out of this thing?
- [ ] Optimize the get transactions call to use a single query; full join, order and then iterate to make transactions object
- [ ] Add a stress testing system
- [ ] set up dev appvault and set secrets
    - [ ] add default limit size as a configurable somewhere
    - [ ] set up configuration loader or a new secrets loader to load values from appvault
- [ ] redis datastore? Maybe later.
- [x] retry startup dependecies
    - [x] wont do; i would prefer to rely on some container/pod scheduler like k8s
- [x] Events
    - [x] NATS
        - [ ] Run NATS locally in a cluster
    - [x] Instrument account.created, transaction.created events
- [x] improve startup configurations
    - [x] update env config to include configurations for nats
    - [s] add configuration for starting nats wire tap in main-app
- [x] add gRPC methods, matching http spec
    - [x] add gRPC server start to main?
    - [x] add gRPC method tests..? idk
- [x] better startup logging
    - [x] info and such about what's starting.
- [x] make basis an enum
- [x] re-evaluate type validation logic in the service
- [x] implement db types and transformers
    - [x] add tests
- [x] [implement safer enums](https://threedots.tech/post/safer-enums-in-go/)
- [x] migrate to guid from uuid
    - [x] probably requires db migration/update
- [x] migrate to ksuid from uuid
- [x] version the api's
- [x] use httpapi models for request/response of all resources
- [x] decouple application runtime, environment, controller, and test!
- [x] implement golang-migrate or similar db migration strategy
    - [x] include the sql files as bin data in binary so migrator can run them ez pz
- [x] Complete the basic CRUD for book-keeping
- [x] make DB methods for CRUD
- [x] set up linting
- [x] Add metadata fields for all tables (created_at, updated_at)
- [x] Better error handling
    - [x] Better error logging
- [x] Pagination
- [x] swap pq for pgx
- [x] Add code-coverage report via github actions?
    - [x] https://blog.seriesci.com/how-to-measure-code-coverage-in-go/
- [x] Add service layer 
- [x] Decouple app setup and routing
- [x] implement a separate errors package like unto nate finches error flags solution
- [x] split api encode and write code
- [x] swap to uuid over auto incr values for id's
- [x] swap go-chi router logger to use zap logger stored in environment; log more.
- [x] ~~add delete request response with message that this ledger is append only and records cannot be deleted~~ 404 is fine
- [x] fix timestamp for created_at to be utc time zone
- [x] change limit and cursor to optional values
    - [x] use [stable pagination](http://morningcoffee.io/stable-pagination.html) for uuid
- [x] swap to nubanks balanced by design transaction model
    - [x] replace transaction with single value and debit/credit accounts; balanced by design
- [x] change the amount in oas to string, and change amount values to int64 in controller/service
    - [x] because postgres (and most dbs) don't OOB implement unsigned ints, use an int 64, which should be more than enough for any needs we'll have. In the case where a string request comes in (upper unbounded), split into multiple entries which will fit into int63's. Probably implement some overflow checks as well.
- [x] generate zero-dep html file for api docs and create serve w/ makefile commands
- [x] implement standard api response and error response to simplify controll handler functions
    - [x] more elegant error response handling
    - [x] better logging
- [x] ~~Use https://mermade.github.io/openapi-gui/ to generate OAS and serve it~~ I went with hand-rolled .yml and redocly serving
- [x] Give config.Database a method to create a connection string
    - [x] Replace NewDatabase with a single conn string vs broken out vars?
- [x] Improve convention for environment variable key constants?
- [x] set environment types as consts

