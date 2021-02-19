# GoLuca

A Simple Accounting Ledger

[![Go Report Card](https://goreportcard.com/badge/github.com/abelgoodwin1988/GoLuca)](https://goreportcard.com/report/github.com/abelgoodwin1988/GoLuca) [![Coverage Status](https://coveralls.io/repos/github/abelgoodwin1988/GoLuca/badge.svg?branch=main)](https://coveralls.io/github/abelgoodwin1988/GoLuca?branch=main)

- Simple application which writes and reads accounting ledger entries to a postgres database

Database

- Postgres database with a single Schema and single Entries table
- For development we'll be blowing up the database every time, and creating the schema (seeding soon TM)
- Later we'll introduce a migration proess for my self-learning

Secrets

- Secrets will be read from environment variables

Making the environments ready with data

- We'll be using [mage](https://magefile.org/) to preload a development appvault with secrets

TODO

- [x] Complete the basic CRUD for book-keeping
- [x] make DB methods for CRUD
- [x] set up linting
- [x] Add metadata fields for all tables (created_at, updated_at)
- [x] Better error handling
    - [x] Better error logging
- [x] Pagination
- [ ] Add a queueing system?
    - idk about this, we need transactions to be completely ACID.. maybe check queue if no affected accounts...?
- [ ] Add some kind of concurrency for the luls and the learnings
- [ ] Add a stress testing system
    - [ ] Magefile
- [ ] Decouple app setup and routing
    - [ ] magefile for chi doc generation
- [ ] Add swaggo b/c I'm lazy at curl
    - [ ] https://github.com/swaggo/swag
- [ ] Create a seeder for a basic dev environment of data
- [x] Add code-coverage report via github actions?
    - [x] https://blog.seriesci.com/how-to-measure-code-coverage-in-go/

