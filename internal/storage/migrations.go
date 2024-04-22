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
	`CREATE TABLE IF NOT EXISTS migration_test_table (
		id text PRIMARY KEY,
		name text,
		logo text,
		website text
	);`,
	`CREATE TABLE IF NOT EXISTS maker_pie_yearly
  (
      year int,
      makerid text,
      id text,
      displayname text,
      fresh boolean,
      labels text[],
      image_file text,
      web_link text,
      pack_count int,
      pack_price_in_pence int,
      PRIMARY KEY (year, makerId, id)
  );`,
	`CREATE TABLE IF NOT EXISTS maker_pie_ranking_yearly
  (
      year int,
      makerid text,
      pieid text,
      userid text,
      pastry int,
      filling int,
      topping int,
      looks int,
      value int,
      notes text,
      last_updated timestamp,
      PRIMARY KEY (year, makerid, pieid, userid)
  );`,
	`CREATE TABLE IF NOT EXISTS config (
	key text PRIMARY KEY,
	value text
  )`,
	`CREATE TABLE IF NOT EXISTS admins (
	id text PRIMARY KEY
  )`,
	`CREATE TABLE IF NOT EXISTS user_pie_yearly
  (
      year int,
      id text,
      owner_userid text,
      maker text,
      location text,
      displayname text,
      fresh boolean,
      labels text[],
      image_file text,
      web_link text,
      pack_count int,
      pack_price_in_pence int,
      clean boolean,
      PRIMARY KEY (year, id, owner_userid)
  );`,
	`CREATE TABLE IF NOT EXISTS user_pie_ranking_yearly
  (
      year int,
      pieid text,
      userid text,
      pastry int,
      filling int,
      topping int,
      looks int,
      value int,
      notes text,
      last_updated timestamp,
      PRIMARY KEY (year, pieid, userid)
  );`,
	`CREATE TABLE IF NOT EXISTS maker (
		id text PRIMARY KEY,
		name text,
		logo text,
		website text
	)`,
	`ALTER TABLE maker_pie_yearly ADD validated boolean;`,
	`CREATE TABLE IF NOT EXISTS categories (
		id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
		slug text,
		label text
	)`,
	`ALTER TABLE maker_pie_yearly ADD oid uuid DEFAULT gen_random_uuid();`,
	`CREATE TABLE IF NOT EXISTS maker_pie_categories (
		maker_pie_oid uuid,
		category_id uuid,
		PRIMARY KEY (maker_pie_oid, category_id)
	)`,
	`WITH all_categories as (
	  SELECT DISTINCT(unnest(mpy.labels)) as category
	  FROM maker_pie_yearly mpy
	), slugged_categories as (
			SELECT REPLACE(category, ' ', '-') slug
	  FROM all_categories
	), distinct_slugs as (
	  SELECT DISTINCT slug
	  FROM slugged_categories
	)
	INSERT INTO categories
	select gen_random_uuid() as id, slug,
			INITCAP(REPLACE(slug, '-', ' ')) as label
	from distinct_slugs`,
	`with all_categories as (
	  SELECT DISTINCT(unnest(mpy.labels)) as category, mpy.oid
	  FROM maker_pie_yearly mpy
	), clean_categories as (
			SELECT REPLACE(category, ' ', '-') as slug, oid
	  FROM all_categories
	)
	INSERT INTO maker_pie_categories
	select mpy.oid as maker_pie_oid,
			c.id as category_id
	FROM maker_pie_yearly mpy
	inner join clean_categories cc on cc.oid = mpy.oid
	inner join categories c on c.slug=cc.slug`,
	`
	CREATE VIEW maker_pie_yearly_categories AS
		  SELECT
		    mpy.oid as oid,
		    jsonb_agg(jsonb_build_object('id', c.id, 'slug', c.slug, 'label', c.label)) as categories,
			ARRAY_AGG(c.slug) as category_slugs
		  from maker_pie_yearly mpy
		  inner join maker_pie_categories mpc on mpc.maker_pie_oid = mpy.oid
		  inner join categories c on mpc.category_id = c.id
		  group by mpy.oid;
	`,
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
