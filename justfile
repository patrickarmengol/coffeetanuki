set dotenv-load

# run app server
run:
    go run ./cmd/web -db-dsn=${CT_DB_DSN}

# access db with psql
psql:
    psql ${CT_DB_DSN}

