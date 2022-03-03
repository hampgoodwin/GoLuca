# Migrations

We use [golang-migrate/migrate](https://github.com/golang-migrate/migrate).

A local install of the golang-migrate/migreate CLI tool is not expected, nor run it as such. We use the tool as a library. We embed the migrations directory in the go binary, along with the golang-migrate/migrate library package. The package then uses the embeded files to run migrations.

So, when doing local development, the migration of the database resource is tied to the app. If you wish to run migrations independent of the application, the CLI tool is required.

## Migration Files

We maintain migration files as .sql in the database package directory. We keep it here to make it easier to embed the sql files via go:embed.

My experience is that a migration succeeds or fails gracefully. For this reason, we do not provide migration down files. In the event of a mistake, roll forward with a new migration. In the event of a dirty database, roll forward.
