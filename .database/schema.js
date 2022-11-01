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
      displayName text,
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
      pastry int,
      filling int,
      topping int,
      looks int,
      value int,
      notes text
  );`,
  `CREATE INDEX IF NOT EXISTS ON mincepierank.maker_pie_ranking (userId);`,
  `CREATE INDEX IF NOT EXISTS ON mincepierank.maker_pie_ranking (pastry);`,
  `CREATE INDEX IF NOT EXISTS ON mincepierank.maker_pie_ranking (filling);`,
  `CREATE INDEX IF NOT EXISTS ON mincepierank.maker_pie_ranking (topping);`,
  `CREATE INDEX IF NOT EXISTS ON mincepierank.maker_pie_ranking (looks);`,
  `CREATE INDEX IF NOT EXISTS ON mincepierank.maker_pie_ranking (value);`,
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
