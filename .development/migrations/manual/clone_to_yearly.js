const cassandra = require('cassandra-driver');

// This is for local only - no env

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

function rowToObject(row) {
  const constructed = {};
  row.keys().forEach((key) => (constructed[key] = row.get(key)));
  return constructed;
}

const client = new cassandra.Client(MIGRATION_CLIENT_CONFIG);

client.connect().then(async () => {
  const result = await client.execute('SELECT * FROM mincepierank.maker_pie;');

  const makers = result.rows.map(rowToObject);

  await Promise.all(
    makers.map((maker) =>
      client.execute(
        `INSERT INTO 
      mincepierank.maker_pie_yearly (
        year,
        makerid,
        id,
        displayname,
        fresh,
        labels,
        image_file,
        web_link,
        pack_count,
        pack_price_in_pence
      )
      VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`,
        [
          2022,
          maker.makerid,
          maker.id,
          maker.displayname,
          maker.fresh,
          maker.labels,
          maker.image_file,
          maker.web_link,
          maker.pack_count,
          maker.pack_price_in_pence
        ],
        { prepare: true }
      )
    )
  );
  const ranking_result = await client.execute('SELECT * FROM mincepierank.maker_pie_ranking;');

  const rankings = ranking_result.rows.map(rowToObject);

  await Promise.all(
    rankings.map((ranking) =>
      client.execute(
        `INSERT INTO 
        mincepierank.maker_pie_ranking_yearly (
          year,
          makerid,
          pieid,
          userid,
          pastry,
          filling,
          topping,
          looks,
          value,
          notes,
          last_updated
        )
        VALUES (
          ?,
          ?,
          ?,
          ?,
          ?,
          ?,
          ?,
          ?,
          ?,
          ?,
          ?
        );`,
        [
          2022,
          ranking.makerid,
          ranking.pieid,
          ranking.userid,
          ranking.pastry,
          ranking.filling,
          ranking.topping,
          ranking.looks,
          ranking.value,
          ranking.notes,
          ranking.last_updated
        ],
        { prepare: true }
      )
    )
  );
  process.exit(0);
});
