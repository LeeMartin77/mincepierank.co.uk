import { CASSANDRA_CLIENT } from './cassandra';

const MIGRATIONS = [
  [
    '000001',
    `CREATE TABLE IF NOT EXISTS mincepierank.maker (
    id text PRIMARY KEY,
    name text,
    logo text,
    website text
)`
  ],
  [
    '000002',
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
  );`
  ],
  ['000003', `CREATE INDEX IF NOT EXISTS ON mincepierank.maker_pie_yearly (year);`],
  ['000004', `CREATE INDEX IF NOT EXISTS ON mincepierank.maker_pie_yearly (makerId);`],
  ['000005', `CREATE INDEX IF NOT EXISTS ON mincepierank.maker_pie_yearly (fresh);`],
  ['000006', `CREATE INDEX IF NOT EXISTS ON mincepierank.maker_pie_yearly (labels);`],
  [
    '000007',
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
  );`
  ],
  ['000008', `CREATE INDEX IF NOT EXISTS ON mincepierank.maker_pie_ranking_yearly (year);`],
  ['000009', `CREATE INDEX IF NOT EXISTS ON mincepierank.maker_pie_ranking_yearly (makerid);`],
  ['000010', `CREATE INDEX IF NOT EXISTS ON mincepierank.maker_pie_ranking_yearly (pieid);`],
  ['000011', `CREATE INDEX IF NOT EXISTS ON mincepierank.maker_pie_ranking_yearly (userId);`],
  ['000012', `CREATE INDEX IF NOT EXISTS ON mincepierank.maker_pie_ranking_yearly (pastry);`],
  ['000013', `CREATE INDEX IF NOT EXISTS ON mincepierank.maker_pie_ranking_yearly (filling);`],
  ['000014', `CREATE INDEX IF NOT EXISTS ON mincepierank.maker_pie_ranking_yearly (topping);`],
  ['000015', `CREATE INDEX IF NOT EXISTS ON mincepierank.maker_pie_ranking_yearly (looks);`],
  ['000016', `CREATE INDEX IF NOT EXISTS ON mincepierank.maker_pie_ranking_yearly (value);`],
  ['000017', `CREATE INDEX IF NOT EXISTS ON mincepierank.maker_pie_ranking_yearly (last_updated);`]
];

const SEEDS = [];

export const runMigrations = async () => {
  await CASSANDRA_CLIENT.execute(
    `CREATE KEYSPACE IF NOT EXISTS mincepierank WITH REPLICATION = {'class':'SimpleStrategy','replication_factor':1};`
  );
  await CASSANDRA_CLIENT.execute(
    `CREATE TABLE IF NOT EXISTS mincepierank.migrations (key text PRIMARY KEY)`
  );
  for (const [key, migration] of MIGRATIONS) {
    const res = await CASSANDRA_CLIENT.execute(
      `SELECT key FROM mincepierank.migrations WHERE key = ?`,
      [key],
      { prepare: true }
    );
    if (res.rows.length === 0) {
      console.log('Running migration: ', key);
      await CASSANDRA_CLIENT.execute(migration);
      await CASSANDRA_CLIENT.execute(
        `INSERT INTO mincepierank.migrations (key) VALUES (?)`,
        [key],
        {
          prepare: true
        }
      );
    }
  }
  for (const [seedCommand, seedEntries] of SEEDS) {
    for (const seedEntry of seedEntries)
      await CASSANDRA_CLIENT.execute(
        // eslint-disable-next-line @typescript-eslint/ban-ts-comment
        // @ts-ignore
        // eslint-disable-next-line @typescript-eslint/no-unused-vars
        `${seedCommand} values (${seedEntry.map((_) => '?').join(',')});`,
        seedEntry,
        { prepare: true }
      );
  }
};
