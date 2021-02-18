Application

- Simple application which writes and reads accounting ledger entries to a postgres database

Database

- Postgres database with a single Schema and single Entries table
- For development we'll be blowing up the database every time, and creating the schema
- Later we'll introduce a migration proess for my self-learning

Secrets

- Secrets will be read from environment variables

Making the environments ready with data

- We'll be using [mage](https://magefile.org/) to preload a development appvault with secrets

TODO

- [x] Complete the basic CRUD for book-keeping
- [x] make DB methods for CRUD
- [ ] set up linting
- [ ] Add metadata fields for all tables (created_at, updated_at)
- [ ] Better error handling
    - [ ] Better error logging
- [ ] Pagination
- [ ] Cache for account name/data?
- [ ] Add a queueing system
    - [ ] Add some kind of concurrency for the luls and learnings
- [ ] Add a stress testing system

