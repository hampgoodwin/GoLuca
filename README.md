# GoLuca

A Simple Accounting Ledger

[![Go Report Card](https://goreportcard.com/badge/github.com/hampgoodwin/GoLuca)](https://goreportcard.com/report/github.com/hampgoodwin/GoLuca) [![Coverage Status](https://coveralls.io/repos/github/hampgoodwin/GoLuca/badge.svg)](https://coveralls.io/github/hampgoodwin/GoLuca) [![golangci-lint](https://github.com/hampgoodwin/GoLuca/actions/workflows/golint-ci.yml/badge.svg)](https://github.com/hampgoodwin/GoLuca/actions/workflows/golint-ci.yml)

- Simple application which writes and reads accounting ledger entries

## Tooling

- [gofumpt](https://github.com/mvdan/gofumpt) for formatting.
- [golangci-lint](https://github.com/golangci/golangci-lint) for linting.
- [colima](https://github.com/abiosoft/colima) for container runtimes.

---

TODO
- [ ] implement db types and transformers
- [ ] add gRPC methods, matching http spec
- [ ] Testing
    - [x] Add more unit tests, solitary and sociable
    - [x] Add Integration testing
    - [x] Add fail case account controller testing
    - [x] Use httptest for tests instead of setting up an http server for each test
        - [x] Fix bug where connection is rejected (I sus due to db+http connection limitations in unix)
    - [ ] Add cursor testing for listing accounts
    - [ ] Add fuzzing
    - [ ] Create a seeder for a basic dev environment of data
    - [ ] Do some stress testing; how much data and we throw at & get out of this thing?
- [ ] Optimize the get transactions call to use a single query; full join, order and then iterate to make transactions object
- [ ] Add a stress testing system
- [ ] set up dev appvault and set secrets
    - [ ] add default limit size as a configurable somewhere
- [ ] improve pagination by displaying page number of result.
    - [ ] set up configuration loader or a new secrets loader to load values from appvault
- [ ] make a frontend with some dashboard functionality, (vue3 plz) OOOOH!!!
- [ ] Events
    - [ ] Add eventing
    - [ ] explore NATS
- [ ] o11y
    - [ ] Add tracing
    - [ ] Add metrics
        - [ ] OpenTelemetry metrics ready?
    - [ ] Log Collector?
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

