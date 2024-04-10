const pg = require("pg");
const cassandra = require("cassandra-driver");
const CassClient = cassandra.Client;
const Mapper = cassandra.mapping.Mapper;

const { Client } = pg;
const connectionString =
  process.env.POSTGRESS_CONNECTION_STRING ?? "postgresql://postgres:mysecretpassword@localhost:5432/postgres";
const cassClient = new CassClient({
  contactPoints: process.env.CASSANDRA_HOST ? process.env.CASSANDRA_HOST.split(";") :  ["localhost:9143"],
  localDataCenter: process.env.CASSANDRA_LOCALDATACENTER ?? 'datacenter1',
  credentials: {
    username: process.env.CASSANDRA_USER ?? 'cassandra',
    password: process.env.CASSANDRA_PASSWORD ?? 'cassandra'
  },
  keyspace: "mincepierank",
});
const tablestoexport = [
  "admins",
  "config",
  "maker_pie_ranking_yearly",
  "maker_pie_yearly",
  "maker",
  "user_pie_ranking_yearly",
  "user_pie_yearly",
];
const client = new Client({
  connectionString,
});
async function main() {
  await client.connect();
  for (const table of tablestoexport) {
    console.log("migrating ", table);
    const mapper = new Mapper(cassClient, {
      models: { [table]: { tables: [table] } },
    });
    const userPieMapper = mapper.forModel(table);
    const items = await userPieMapper.findAll();
    for (let item of items) {
      const query = {
        text: `INSERT INTO ${table} (${Object.keys(item).join(
          ", "
        )}) VALUES (${Object.keys(item)
          .map((x, i) => `$${i + 1}`)
          .join(", ")}) ON CONFLICT DO NOTHING`,
        values: Object.values(item),
      };

      await client.query(query);
    }
  }

  console.log("migrated successfully");
  process.exit(0);
}

main();
