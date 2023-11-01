set dotenv-load

# run app server
run:
    go run ./cmd/web -db-dsn=${CT_DB_DSN}

# access db with psql
psql:
    psql ${CT_DB_DSN}

# create a new migration
migratenew name:
    migrate create -seq -ext=.sql -dir=./migrations {{name}}

# migrate up all the way
migrateup:
    migrate -path=./migrations -database=${CT_DB_DSN} up

# migrate down all the way
migratedown:
    migrate -path=./migrations -database=${CT_DB_DSN} down

# migrate to a specific version
migrateto version:
    migrate -path=./migrations -database=${CT_DB_DSN} goto {{version}}

# force the migration version number
migrateforce version:
    migrate -path=./migrations -database=${CT_DB_DSN} force {{version}}
