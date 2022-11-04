const cassandra = require("cassandra-driver");
const { randomUUID } = require("crypto");

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
  {
    makerId: "tesco",
    id: "not-nice",
    displayName: "Okay Nice Pie",
    fresh: false,
    labels: ["okay", "mediocre"],
  },
  {
    makerId: "morrisons",
    id: "super-nice",
    displayName: "Super Good Pie",
    fresh: false,
    labels: ["amazing"],
  },
  {
    makerId: "morrisons",
    id: "nice",
    displayName: "Okay Pie",
    fresh: false,
    labels: ["okay"],
  },
];

const pieRankings = makerPies.reduce((rankings, pie) => {
  const newRankings = [1, 2, 3, 4, 5].map((number) => {
    return {
      makerid: pie.makerId,
      pieid: pie.id,
      userid: randomUUID(),
      pastry: number,
      filling: number,
      topping: number,
      looks: number,
      value: number,
      notes: "Some Notes for " + number + " :: " + randomUUID(),
    };
  });
  return [...rankings, ...newRankings];
}, []);

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

  for (let pieRanking of pieRankings) {
    await client.execute(
      `INSERT INTO mincepierank.maker_pie_ranking (makerid, pieid, userid, pastry, filling, topping, looks, value, notes)
            values (?, ?, ?, ?, ?, ?, ?, ?, ?);`,
      Object.values(pieRanking),
      { prepare: true }
    );
  }
  return client.shutdown();
});
