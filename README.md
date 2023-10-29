# coffeetanuki

Website for rating/reviewing coffee beans and roasters.

## Installation

1. Set up PostgreSQL DB.
2. Place DSN in a `.env` file under variable `CT_DB_DSN`.

## Usage

> This project uses [just](https://github.com/casey/just) for command running.

```
$ just -l
Available recipes:
    migratedown # migrate down
    migrateup   # migrate up
    psql        # access db with psql
    run         # run app server
```

## License

The source code for `coffeetanuki` is distributed under the terms of any of the following licenses:

- [Apache-2.0](https://spdx.org/licenses/Apache-2.0.html)
- [MIT](https://spdx.org/licenses/MIT.html)

The content on the site is distributed under:

- [CC BY NC SA 4.0](https://spdx.org/licenses/CC-BY-NC-SA-4.0.html)
