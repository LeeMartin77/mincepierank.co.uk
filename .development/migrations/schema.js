const cassandra = require('cassandra-driver');

const schema = [
	`CREATE KEYSPACE IF NOT EXISTS mincepierank 
  WITH REPLICATION = {'class':'SimpleStrategy','replication_factor':1};`,
	`CREATE TABLE IF NOT EXISTS mincepierank.maker (
      id text PRIMARY KEY,
      name text,
      logo text,
      website text
  )`,
	`CREATE TABLE IF NOT EXISTS mincepierank.maker_pie_yearly
  (
      year int,
      makerid text,
      id text,
      displayname text,
      fresh boolean,
      labels set<text>,
      image_file text,
      web_link text,
      pack_count int,
      pack_price_in_pence int,
      PRIMARY KEY ((year, makerId, id))
  );`,
	`CREATE INDEX IF NOT EXISTS ON mincepierank.maker_pie_yearly (year);`,
	`CREATE INDEX IF NOT EXISTS ON mincepierank.maker_pie_yearly (makerId);`,
	`CREATE INDEX IF NOT EXISTS ON mincepierank.maker_pie_yearly (fresh);`,
	`CREATE INDEX IF NOT EXISTS ON mincepierank.maker_pie_yearly (labels);`,
	`CREATE TABLE IF NOT EXISTS mincepierank.maker_pie_ranking_yearly
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
      PRIMARY KEY ((year, makerid, pieid), userid)
  );`,
	`CREATE INDEX IF NOT EXISTS ON mincepierank.maker_pie_ranking_yearly (year);`,
	`CREATE INDEX IF NOT EXISTS ON mincepierank.maker_pie_ranking_yearly (makerid);`,
	`CREATE INDEX IF NOT EXISTS ON mincepierank.maker_pie_ranking_yearly (pieid);`,
	`CREATE INDEX IF NOT EXISTS ON mincepierank.maker_pie_ranking_yearly (userId);`,
	`CREATE INDEX IF NOT EXISTS ON mincepierank.maker_pie_ranking_yearly (pastry);`,
	`CREATE INDEX IF NOT EXISTS ON mincepierank.maker_pie_ranking_yearly (filling);`,
	`CREATE INDEX IF NOT EXISTS ON mincepierank.maker_pie_ranking_yearly (topping);`,
	`CREATE INDEX IF NOT EXISTS ON mincepierank.maker_pie_ranking_yearly (looks);`,
	`CREATE INDEX IF NOT EXISTS ON mincepierank.maker_pie_ranking_yearly (value);`,
	`CREATE INDEX IF NOT EXISTS ON mincepierank.maker_pie_ranking_yearly (last_updated);`
];

const MIGRATION_CLIENT_CONFIG = {
	contactPoints: process.env.CASSANDRA_CONTACT_POINTS
		? process.env.CASSANDRA_CONTACT_POINTS.split(';')
		: ['localhost:9143'],
	localDataCenter: process.env.CASSANDRA_LOCALDATACENTER ?? 'datacenter1',
	credentials: {
		username: process.env.CASSANDRA_USER ?? 'cassandra',
		password: process.env.CASSANDRA_PASSWORD ?? 'cassandra'
	}
};

const client = new cassandra.Client(MIGRATION_CLIENT_CONFIG);

client.connect().then(async () => {
	for (let migration of schema) {
		await client.execute(migration);
	}
	await client.shutdown();
	process.exit(0);
});
