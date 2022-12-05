const cassandra = require("cassandra-driver");

const schema = [
  `CREATE KEYSPACE IF NOT EXISTS mincepierank 
  WITH REPLICATION = {'class':'SimpleStrategy','replication_factor':1};`,
  `CREATE TABLE IF NOT EXISTS mincepierank.maker (
      id text PRIMARY KEY,
      name text,
      logo text,
      website text
  )`,
  `CREATE TABLE IF NOT EXISTS mincepierank.maker_pie
  (
      makerid text,
      id text,
      displayname text,
      fresh boolean,
      labels set<text>,
      image_file text,
      web_link text,
      pack_count int,
      pack_price_in_pence int,
      PRIMARY KEY ((makerId, id))
  );`,
  `CREATE INDEX IF NOT EXISTS ON mincepierank.maker_pie (makerId);`,
  `CREATE INDEX IF NOT EXISTS ON mincepierank.maker_pie (fresh);`,
  `CREATE INDEX IF NOT EXISTS ON mincepierank.maker_pie (labels);`,
  `CREATE TABLE IF NOT EXISTS mincepierank.maker_pie_ranking
  (
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
      PRIMARY KEY ((makerid, pieid), userid)
  );`,
  `CREATE INDEX IF NOT EXISTS ON mincepierank.maker_pie_ranking (makerid);`,
  `CREATE INDEX IF NOT EXISTS ON mincepierank.maker_pie_ranking (pieid);`,
  `CREATE INDEX IF NOT EXISTS ON mincepierank.maker_pie_ranking (userId);`,
  `CREATE INDEX IF NOT EXISTS ON mincepierank.maker_pie_ranking (pastry);`,
  `CREATE INDEX IF NOT EXISTS ON mincepierank.maker_pie_ranking (filling);`,
  `CREATE INDEX IF NOT EXISTS ON mincepierank.maker_pie_ranking (topping);`,
  `CREATE INDEX IF NOT EXISTS ON mincepierank.maker_pie_ranking (looks);`,
  `CREATE INDEX IF NOT EXISTS ON mincepierank.maker_pie_ranking (value);`,
  `CREATE INDEX IF NOT EXISTS ON mincepierank.maker_pie_ranking (last_updated);`,
];

const MIGRATION_CLIENT_CONFIG = {
  contactPoints: process.env.CASSANDRA_CONTACT_POINTS
    ? process.env.CASSANDRA_CONTACT_POINTS.split(";")
    : ["localhost:9143"],
  localDataCenter: process.env.CASSANDRA_LOCALDATACENTER ?? "datacenter1",
  credentials: {
    username: process.env.CASSANDRA_USER ?? "cassandra",
    password: process.env.CASSANDRA_PASSWORD ?? "cassandra",
  },
};

const client = new cassandra.Client(MIGRATION_CLIENT_CONFIG);

client.connect().then(async () => {
  for (let migration of schema) {
    await client.execute(migration);
  }
  return client.shutdown();
});
