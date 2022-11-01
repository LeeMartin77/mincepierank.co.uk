const cassandra = require("cassandra-driver");

const schema = [
  `CREATE KEYSPACE IF NOT EXISTS mincepierank 
  WITH REPLICATION = {'class':'SimpleStrategy','replication_factor':1};`,
  `CREATE TABLE IF NOT EXISTS mincepierank.maker (
      id text PRIMARY KEY,
      name text
  )`,
  `CREATE TABLE IF NOT EXISTS mincepierank.maker_pie
  (
      makerId text,
      id text,
      display_name text,
      fresh boolean,
      labels set<text>,
      PRIMARY KEY ((makerId, id))
  );`,
  `CREATE INDEX IF NOT EXISTS ON mincepierank.maker_pie (fresh);`,
  `CREATE INDEX IF NOT EXISTS ON mincepierank.maker_pie (labels);`,
  `CREATE TABLE IF NOT EXISTS mincepierank.maker_pie_ranking
(
    pieId text PRIMARY KEY,
    userId text,
    metrics set<int>,
    notes text
);`,
  `CREATE INDEX IF NOT EXISTS ON mincepierank.maker_pie_ranking (userId);`,
  `CREATE INDEX IF NOT EXISTS ON mincepierank.maker_pie_ranking (metrics);`,
];

const MIGRATION_CLIENT_CONFIG = {
  contactPoints: process.env.OPENFOODDIARY_CASSANDRA_CONTACT_POINTS
    ? process.env.OPENFOODDIARY_CASSANDRA_CONTACT_POINTS.split(";")
    : ["localhost:9143"],
  localDataCenter:
    process.env.OPENFOODDIARY_CASSANDRA_LOCALDATACENTER ?? "datacenter1",
  credentials: {
    username: process.env.OPENFOODDIARY_CASSANDRA_USER ?? "cassandra",
    password: process.env.OPENFOODDIARY_CASSANDRA_PASSWORD ?? "cassandra",
  },
};

const client = new cassandra.Client(MIGRATION_CLIENT_CONFIG);

client.connect().then(async () => {
  for (let migration of schema) {
    await client.execute(migration);
  }
  return client.shutdown();
});
