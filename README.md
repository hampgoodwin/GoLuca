# GoLuca

A Simple Accounting Ledger

[![Go Report Card](https://goreportcard.com/badge/github.com/hampgoodwin/GoLuca)](https://goreportcard.com/report/github.com/hampgoodwin/GoLuca) [![Coverage Status](https://coveralls.io/repos/github/hampgoodwin/GoLuca/badge.svg?branch=main)](https://coveralls.io/github/hampgoodwin/GoLuca?branch=main) [![golangci-lint](https://github.com/hampgoodwin/GoLuca/actions/workflows/golint-ci.yml/badge.svg)](https://github.com/hampgoodwin/GoLuca/actions/workflows/golint-ci.yml)

- Simple application which writes and reads accounting ledger entries

TODO

- [ ] change limit and cursor to optional values
    - [x] use [stable pagination](http://morningcoffee.io/stable-pagination.html) for uuid
- [ ] Optimize the get transactions call to use a single query; full join, order and then iterate to make transactions object
- [ ] swap to nubanks balanced by design transaction model
    - [ ] replace transaction with single value and debit/credit accounts; balanced by design
    - [ ] change the amount in oas to string, and change amount values to uint64
- [ ] implement standard api response and error response to simplify api handler functions
    - [ ] more elegant error response handling
- [ ] implement golang-migrate or similar db migration strategy
- [ ] Use https://mermade.github.io/openapi-gui/ to generate OAS and serve it
- [ ] Add a stress testing system
- [ ] Add fuzzing
- [ ] Create a seeder for a basic dev environment of data
- [ ] set up dev appvault and set secrets
- [ ] improve pagination by displaying page number of result.
    - [ ] set up configuration loader or a new secrets loader to load values from appvault
- [ ] make a frontend with some dashboard functionality, (vue3 plz) OOOOH!!!
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

