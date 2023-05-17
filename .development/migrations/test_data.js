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
  { 
    id: "tesco",
    name: "Tesco",
    logo: "https://static.mincepierank.co.uk/images/tesco/logo.png", 
    website: "https://www.tesco.com/"
  },
  { 
    id: "morrisons",
    name: "Morrisons",
    logo: "https://static.mincepierank.co.uk/images/morrisons/logo.png", 
    website: "https://groceries.morrisons.com/webshop/startWebshop.do"
  },
];

const makerPies = [
  {
    makerId: "tesco",
    id: "very-nice",
    displayName: "Very Nice Pie",
    fresh: false,
    labels: ["classic", "vegan"],
    image_file: "https://static.mincepierank.co.uk/images/waitrose/all-butter.jpg",
    web_link: "https://www.waitrose.com/ecom/products/waitrose-christmas-all-butter-mince-pies/735193-785976-785977",
    pack_count: 6,
    pack_price_in_pence: 120
  },
  {
    makerId: "tesco",
    id: "not-nice",
    displayName: "Okay Nice Pie",
    fresh: false,
    labels: ["premium", "mini"],
    image_file: "https://static.mincepierank.co.uk/images/waitrose/all-butter.jpg",
    web_link: "https://www.waitrose.com/ecom/products/waitrose-christmas-all-butter-mince-pies/735193-785976-785977",
    pack_count: 6,
    pack_price_in_pence: 120
  },
  {
    makerId: "morrisons",
    id: "super-nice",
    displayName: "Super Good Pie",
    fresh: false,
    labels: ["alternative", "gluten free"],
    image_file: "https://static.mincepierank.co.uk/images/waitrose/all-butter.jpg",
    web_link: "https://www.waitrose.com/ecom/products/waitrose-christmas-all-butter-mince-pies/735193-785976-785977",
    pack_count: 6,
    pack_price_in_pence: 120
  },
  {
    makerId: "morrisons",
    id: "nice",
    displayName: "Okay Pie",
    fresh: false,
    labels: ["classic", "mini"],
    image_file: "https://static.mincepierank.co.uk/images/waitrose/all-butter.jpg",
    web_link: "https://www.waitrose.com/ecom/products/waitrose-christmas-all-butter-mince-pies/735193-785976-785977",
    pack_count: 6,
    pack_price_in_pence: 120
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
      last_updated: new Date().toISOString()
    };
  });
  return [...rankings, ...newRankings];
}, []);

const client = new cassandra.Client(MIGRATION_CLIENT_CONFIG);

client.connect().then(async () => {
  await client.execute("TRUNCATE mincepierank.maker");
  for (let maker of makers) {
    await client.execute(
      `INSERT INTO mincepierank.maker (id, name, logo, website)
            values (?, ?, ?, ?);`,
      Object.values(maker),
      { prepare: true }
    );
  }
  await client.execute("TRUNCATE mincepierank.maker_pie_yearly");

  for (let makerPie of makerPies) {
    await client.execute(
      `INSERT INTO mincepierank.maker_pie_yearly (year, makerId, id, displayName, fresh, labels, image_file, web_link, pack_count, pack_price_in_pence)
            values (2022, ?, ?, ?, ?, ?, ?, ?, ?, ?);`,
      [...Object.values(makerPie)],
      { prepare: true }
    );
  }

  for (let pieRanking of pieRankings) {
    await client.execute(
      `INSERT INTO mincepierank.maker_pie_ranking_yearly (year, makerid, pieid, userid, pastry, filling, topping, looks, value, notes, last_updated)
            values (2022, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`,
      [...Object.values(pieRanking)],
      { prepare: true }
    );
  }
  return client.shutdown();
});
