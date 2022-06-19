# Migrations

Migrations in the API are handled via [the Darwin library](https://go.pkg.dev/github.com/GuiaBolso/darwin). We keep the individual migration files in [the migrations directory](../database/queriers/postgres/migrations/). They get embedded into the end Go binary by the compiler, which would probably get hairy if the size of the folder weren't currently in the ~14kb range.

We employ the strategy of attempting to migrate the database when the server starts, every time. Darwin is smart enough to know when no migrations need to run, so the 99% case (where deploys are made with no schema changes) is very fast, migration is basically a single-query operation to check on the state of migrations.

## Adding a new table

Tables require more steps to add properly than simple schema updates. Be sure to:

- [ ] Update [the database initializer tool](../../cmd/tools/db_initializer/db_initializer.go) so that it deletes the new table when it first runs.
