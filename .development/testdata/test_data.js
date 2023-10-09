import cassandra from 'cassandra-driver';
import { randomUUID } from 'crypto';

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

const makers = [
  {
    id: 'tesco',
    name: 'Tesco',
    logo: '/tesco-logo.png',
    website: 'https://www.tesco.com/'
  },
  {
    id: 'morrisons',
    name: 'Morrisons',
    logo: '/morrisons-logo.png',
    website: 'https://groceries.morrisons.com/webshop/startWebshop.do'
  }
];

const makerPies = [
  {
    makerId: 'tesco',
    id: 'very-nice',
    displayName: 'Very Nice Pie',
    fresh: false,
    labels: ['classic', 'vegan'],
    image_file: '/w-all-butter.jpg',
    web_link:
      'https://www.waitrose.com/ecom/products/waitrose-christmas-all-butter-mince-pies/735193-785976-785977',
    pack_count: 6,
    pack_price_in_pence: 120,
    validated: true
  },
  {
    makerId: 'tesco',
    id: 'not-nice',
    displayName: 'Okay Nice Pie',
    fresh: false,
    labels: ['premium', 'mini'],
    image_file: '/w-all-butter.jpg',
    web_link:
      'https://www.waitrose.com/ecom/products/waitrose-christmas-all-butter-mince-pies/735193-785976-785977',
    pack_count: 6,
    pack_price_in_pence: 120,
    validated: true
  },
  {
    makerId: 'morrisons',
    id: 'super-nice',
    displayName: 'Super Good Pie',
    fresh: false,
    labels: ['alternative', 'gluten-free'],
    image_file: '/w-all-butter.jpg',
    web_link:
      'https://www.waitrose.com/ecom/products/waitrose-christmas-all-butter-mince-pies/735193-785976-785977',
    pack_count: 6,
    pack_price_in_pence: 120,
    validated: true
  },
  {
    makerId: 'morrisons',
    id: 'nice',
    displayName: 'Okay Pie',
    fresh: false,
    labels: ['classic', 'mini'],
    image_file: '/w-all-butter.jpg',
    web_link:
      'https://www.waitrose.com/ecom/products/waitrose-christmas-all-butter-mince-pies/735193-785976-785977',
    pack_count: 6,
    pack_price_in_pence: 120,
    validated: true
  }
];

const testUserPieRankings = [
  {
    makerid: makerPies[0].makerId,
    pieid: makerPies[0].id,
    userid: 'TEST_AUTH_john.doe@example.com',
    pastry: 5,
    filling: 4,
    topping: 3,
    looks: 2,
    value: 1,
    notes: 'Some Notes for ' + 3 + ' :: ' + randomUUID(),
    last_updated: new Date().toISOString()
  },
  {
    makerid: makerPies[1].makerId,
    pieid: makerPies[1].id,
    userid: 'TEST_AUTH_john.doe@example.com',
    pastry: 4,
    filling: 4,
    topping: 3,
    looks: 5,
    value: 2,
    notes: 'Some Notes for ' + 3 + ' :: ' + randomUUID(),
    last_updated: new Date().toISOString()
  },
  {
    makerid: makerPies[2].makerId,
    pieid: makerPies[3].id,
    userid: 'TEST_AUTH_john.doe@example.com',
    pastry: 1,
    filling: 2,
    topping: 1,
    looks: 2,
    value: 1,
    notes: 'Some Notes for ' + 3 + ' :: ' + randomUUID(),
    last_updated: new Date().toISOString()
  }
];

const pieRankings = [
  ...makerPies.reduce((rankings, pie) => {
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
        notes: 'Some Notes for ' + number + ' :: ' + randomUUID(),
        last_updated: new Date().toISOString()
      };
    });
    return [...rankings, ...newRankings];
  }, []),
  ...testUserPieRankings
];

const client = new cassandra.Client(MIGRATION_CLIENT_CONFIG);

client.connect().then(async () => {
  await client.execute('TRUNCATE mincepierank.maker');
  for (let maker of makers) {
    await client.execute(
      `INSERT INTO mincepierank.maker (id, name, logo, website)
            values (?, ?, ?, ?);`,
      Object.values(maker),
      { prepare: true }
    );
  }
  await client.execute('TRUNCATE mincepierank.maker_pie_yearly');

  for (const year of [2022, 2023]) {
    for (let makerPie of makerPies) {
      await client.execute(
        `INSERT INTO mincepierank.maker_pie_yearly (year, makerId, id, displayName, fresh, labels, image_file, web_link, pack_count, pack_price_in_pence, validated)
              values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`,
        [year, ...Object.values(makerPie)],
        { prepare: true }
      );
    }

    for (let pieRanking of pieRankings) {
      await client.execute(
        `INSERT INTO mincepierank.maker_pie_ranking_yearly (year, makerid, pieid, userid, pastry, filling, topping, looks, value, notes, last_updated)
              values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`,
        [year, ...Object.values(pieRanking)],
        { prepare: true }
      );
    }
  }

  await client.execute(
    `INSERT INTO mincepierank.admins (id) VALUES ('TEST_AUTH_john.doe@example.com')`
  );


  await client.execute(
    `INSERT INTO mincepierank.config (key, value) VALUES ('readonly', 'false')`
  );

  await client.execute(
    `INSERT INTO mincepierank.config (key, value) VALUES ('activeYear', '2023')`
  );
  return client.shutdown();
});
