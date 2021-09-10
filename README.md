# GoLuca

A Simple Accounting Ledger

[![Go Report Card](https://goreportcard.com/badge/github.com/hampgoodwin/GoLuca)](https://goreportcard.com/report/github.com/hampgoodwin/GoLuca) [![Coverage Status](https://coveralls.io/repos/github/hampgoodwin/GoLuca/badge.svg?branch=main)](https://coveralls.io/github/hampgoodwin/GoLuca?branch=main) [![golangci-lint](https://github.com/hampgoodwin/GoLuca/actions/workflows/golint-ci.yml/badge.svg)](https://github.com/hampgoodwin/GoLuca/actions/workflows/golint-ci.yml)

- Simple application which writes and reads accounting ledger entries

TODO

- [x] Complete the basic CRUD for book-keeping
- [x] make DB methods for CRUD
- [x] set up linting
- [x] Add metadata fields for all tables (created_at, updated_at)
- [x] Better error handling
    - [x] Better error logging
- [x] Pagination
- [ ] Add a stress testing system
- [ ] Add fuzzing
- [x] Add service layer 
- [ ] Decouple app setup and routing
- [ ] Use https://mermade.github.io/openapi-gui/ to generate OAS and serve it
- [ ] Create a seeder for a basic dev environment of data
- [x] swap pq for pgx
- [x] Add code-coverage report via github actions?
    - [x] https://blog.seriesci.com/how-to-measure-code-coverage-in-go/
- [ ] set up dev appvault and set secrets
    - [ ] set up configuration loader or a new secrets loader to load values from appvault
- [ ] make a frontend with some dashboard functionality, (vue3 plz) OOOOH!!!
- [ ] improve pagination by displaying page number of result.
- [ ] swap to uuid over auto incr values for id's
- [ ] implement a separate errors package like unto nate finches error flags solution
- [ ] implement standard api response and error response to simplify api handler functions
- [ ] split api encode and write code
- [ ] more elegant error response handling

