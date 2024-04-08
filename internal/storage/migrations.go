package storage

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

var migrations []string = []string{
	`CREATE TABLE IF NOT EXISTS migration_test_table (
		id SERIAL PRIMARY KEY,
		test_data varchar(255),
		created TIMESTAMP NOT NULL DEFAULT NOW()
	);`,
	`INSERT INTO migration_test_table (test_data) VALUES ('test migration data woo')`,
}

var migrationLogTable string = `
	CREATE TABLE IF NOT EXISTS migrations (
		idx integer,
		enabled TIMESTAMP NOT NULL DEFAULT NOW(),
		PRIMARY KEY (idx)
	);
`

var insertLogSql string = `INSERT INTO migrations (idx) VALUES ($1)`

// Migrate implements Operations.
func (o *OperationWrapper) Migrate(c context.Context) error {
	_, err := o.db.Exec(c, migrationLogTable)
	if err != nil {
		log.Error().Err(err).Msg("Error with creating migration log table")
		return err
	}

	for idx, migration := range migrations {
		var sqldx int64
		err := o.db.QueryRow(c, "select idx from migrations where idx=$1", idx).Scan(&sqldx)
		if err == pgx.ErrNoRows {
			// have a migration to do!
			_, err := o.db.Exec(c, migration)
			if err != nil {
				log.Error().Err(err).Msgf("Error with migrating '%d'", idx)
				return err
			}
			_, err = o.db.Exec(c, insertLogSql, idx)
			if err != nil {
				log.Error().Err(err).Msgf("Error with inserting log '%d'", idx)
				return err
			}
		}
	}

	return nil
}
