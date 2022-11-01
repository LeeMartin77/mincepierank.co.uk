const cassandra = require("cassandra-driver");

// This is for local only - no env

const MIGRATION_CLIENT_CONFIG = {
  contactPoints: ["localhost:9143"],
  localDataCenter: "datacenter1",
  credentials: {
    username: "cassandra",
    password: "cassandra",
  },
};

const makers = [
  { id: "tesco", name: "Tesco" },
  { id: "morrisons", name: "Morrisons" },
];

const makerPies = [
  {
    makerId: "tesco",
    id: "very-nice",
    displayName: "Very Nice Pie",
    fresh: false,
    labels: ["good", "value"],
  },
];

const client = new cassandra.Client(MIGRATION_CLIENT_CONFIG);

client.connect().then(async () => {
  await client.execute("TRUNCATE mincepierank.maker");
  for (let maker of makers) {
    await client.execute(
      `INSERT INTO mincepierank.maker (id, name)
            values (?, ?);`,
      Object.values(maker),
      { prepare: true }
    );
  }
  await client.execute("TRUNCATE mincepierank.maker_pie");

  for (let makerPie of makerPies) {
    await client.execute(
      `INSERT INTO mincepierank.maker_pie (makerId, id, displayName, fresh, labels)
            values (?, ?, ?, ?, ?);`,
      Object.values(makerPie),
      { prepare: true }
    );
  }
  return client.shutdown();
});
