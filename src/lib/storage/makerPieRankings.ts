import { err, ok, Result } from 'neverthrow';
import { CASSANDRA_CLIENT } from './cassandra';
import type { MakerPieRanking } from './types';
import { StorageError } from './types';
import { calculateAverage, rowToObject } from './utilities';

export async function getAllRankingsForPie(
  year: number,
  makerid: string,
  pieid: string
): Promise<Result<MakerPieRanking[], StorageError>> {
  try {
    const result = await CASSANDRA_CLIENT.execute(
      'SELECT * FROM mincepierank.maker_pie_ranking_yearly WHERE year = ? AND makerid = ? AND pieid = ? ALLOW FILTERING;',
      [year, makerid, pieid],
      { prepare: true }
    );

    const mapped = result.rows.map(rowToObject);
    return ok(mapped as MakerPieRanking[]);
  } catch {
    return err(StorageError.GenericError);
  }
}

export async function getLatestRanking(): Promise<Result<MakerPieRanking, StorageError>> {
  try {
    const whenresult = await CASSANDRA_CLIENT.execute(
      'SELECT MAX(last_updated) as last_updated FROM mincepierank.maker_pie_ranking_yearly where year = 2022;'
    );

    if (whenresult.rowLength === 0) {
      return err(StorageError.NotFound);
    }

    const when = whenresult.first().get('last_updated');
    if (!when) {
      return err(StorageError.NotFound);
    }

    const top = await CASSANDRA_CLIENT.execute(
      `SELECT year,
      makerid,
      pieid,
      pastry,
      filling,
      topping,
      looks,
      value,
      CAST(last_updated as text) as last_updated FROM mincepierank.maker_pie_ranking_yearly where last_updated >= ? ALLOW FILTERING;`,
      [new Date(when)],
      { prepare: true }
    );
    return ok(rowToObject(top.first()));
  } catch {
    return err(StorageError.GenericError);
  }
}

export async function getMyRankingForPie(
  year: number,
  makerid: string,
  pieid: string,
  userid: string
): Promise<Result<MakerPieRanking, StorageError>> {
  try {
    const result = await CASSANDRA_CLIENT.execute(
      `SELECT year,
      makerid,
      pieid,
      pastry,
      filling,
      topping,
      looks,
      value,
      CAST(last_updated as text) as last_updated FROM mincepierank.maker_pie_ranking_yearly WHERE year = ? makerid = ? AND pieid = ? AND userid = ? ALLOW FILTERING;`,
      [year, makerid, pieid, userid],
      { prepare: true }
    );
    if (result.rows.length === 0) {
      return err(StorageError.NotFound);
    }

    return ok(rowToObject(result.first()) as MakerPieRanking);
  } catch {
    return err(StorageError.GenericError);
  }
}

export async function getUserPieRankings(
  year: number,
  userid: string
): Promise<Result<(MakerPieRanking & { average: number; count: undefined })[], StorageError>> {
  try {
    const result = await CASSANDRA_CLIENT.execute(
      `SELECT year,
        makerid,
        pieid,
        pastry,
        filling,
        topping,
        looks,
        value,
        CAST(last_updated as text) as last_updated
        FROM mincepierank.maker_pie_ranking_yearly WHERE year = ? AND userid = ? ALLOW FILTERING;`,
      [year, userid],
      { prepare: true }
    );

    return ok(
      result.rows.map<MakerPieRanking>(rowToObject).map(addAverageScore) as (MakerPieRanking & {
        count: undefined;
        average: number;
      })[]
    );
  } catch {
    return err(StorageError.GenericError);
  }
}

export type PieRankingSummary = Omit<MakerPieRanking, 'userid' | 'notes'> & {
  count: number | undefined;
  average: number;
};

export function addAverageScore<T>(
  mapped: Pick<MakerPieRanking, 'filling' | 'pastry' | 'topping' | 'value' | 'looks'> & T
) {
  return { ...mapped, average: calculateAverage(mapped) };
}

export async function getPieRankingSummary(
  year: number,
  makerid: string,
  pieid: string
): Promise<Result<PieRankingSummary, StorageError>> {
  try {
    const result = await CASSANDRA_CLIENT.execute(
      `SELECT year,
        makerid,
        pieid,
        AVG(cast(pastry as float)) as pastry,
        AVG(cast(filling as float)) as filling,
        AVG(cast(topping as float)) as topping,
        AVG(cast(looks as float)) as looks,
        AVG(cast(value as float)) as value,
        CAST(MAX(last_updated) as text) as last_updated,
        CAST(COUNT(1) as int) as count 
      FROM mincepierank.maker_pie_ranking_yearly
      WHERE year = ? AND makerid = ? AND pieid = ? ALLOW FILTERING;`,
      [year, makerid, pieid],
      { prepare: true }
    );

    return ok(addAverageScore(rowToObject(result.first())));
  } catch (ex) {
    console.log(ex);
    return err(StorageError.GenericError);
  }
}

export async function getPieRankingSummariesByIds(
  year: number,
  makerid: string,
  pieids: string[]
): Promise<Result<PieRankingSummary[], StorageError>> {
  try {
    const result = await CASSANDRA_CLIENT.execute(
      `SELECT year,
        makerid,
        pieid,
        AVG(cast(pastry as float)) as pastry,
        AVG(cast(filling as float)) as filling,
        AVG(cast(topping as float)) as topping,
        AVG(cast(looks as float)) as looks,
        AVG(cast(value as float)) as value,
        CAST(MAX(last_updated) as text) as last_updated,
        CAST(COUNT(1) as int) as count 
      FROM mincepierank.maker_pie_ranking_yearly
      WHERE year = ? AND makerid = ? and pieid in ?
      GROUP BY makerid, pieid;`,
      [year, makerid, pieids],
      { prepare: true }
    );

    return ok(result.rows.map<PieRankingSummary>(rowToObject).map(addAverageScore));
  } catch {
    return err(StorageError.GenericError);
  }
}

export async function getMakerPieRankingSummaries(
  year: number,
  makerid: string
): Promise<Result<PieRankingSummary[], StorageError>> {
  try {
    const result = await CASSANDRA_CLIENT.execute(
      `SELECT year,
        makerid,
        pieid,
        AVG(cast(pastry as float)) as pastry,
        AVG(cast(filling as float)) as filling,
        AVG(cast(topping as float)) as topping,
        AVG(cast(looks as float)) as looks,
        AVG(cast(value as float)) as value,
        CAST(MAX(last_updated) as text) as last_updated,
        CAST(COUNT(1) as int) as count 
      FROM mincepierank.maker_pie_ranking_yearly
      WHERE year = ? AND makerid = ?
      GROUP BY year, makerid, pieid;`,
      [year, makerid],
      { prepare: true }
    );

    return ok(result.rows.map<PieRankingSummary>(rowToObject).map(addAverageScore));
  } catch {
    return err(StorageError.GenericError);
  }
}

export async function getAllPieRankingSummaries(
  year: number
): Promise<Result<PieRankingSummary[], StorageError>> {
  try {
    const result = await CASSANDRA_CLIENT.execute(
      `SELECT year,
        makerid,
        pieid,
        AVG(cast(pastry as float)) as pastry,
        AVG(cast(filling as float)) as filling,
        AVG(cast(topping as float)) as topping,
        AVG(cast(looks as float)) as looks,
        AVG(cast(value as float)) as value,
        CAST(MAX(last_updated) as text) as last_updated,
        CAST(COUNT(1) as int) as count 
      FROM mincepierank.maker_pie_ranking_yearly
      WHERE year = ?
      GROUP BY year, makerid, pieid;`,
      [year],
      { prepare: true }
    );

    return ok(result.rows.map<PieRankingSummary>(rowToObject).map(addAverageScore));
  } catch {
    return err(StorageError.GenericError);
  }
}

export async function addPieRanking(
  ranking: MakerPieRanking
): Promise<Result<boolean, StorageError>> {
  const mapRanking = ranking as {
    [key: string]: string | number | Date | undefined;
  };
  const rankingValues = ['pastry', 'filling', 'topping', 'looks', 'value'];

  rankingValues.forEach(
    (key) => (mapRanking[key] = parseInt(mapRanking[key]?.toString() ?? 'NaN'))
  );
  if (
    !rankingValues.every(
      (key) => (mapRanking[key] as number) < 6 && (mapRanking[key] as number) > 0
    ) ||
    (ranking.notes && ranking.notes.length > 140) ||
    !ranking.makerid ||
    !ranking.pieid ||
    !ranking.userid
  ) {
    return err(StorageError.BadInput);
  }
  try {
    await CASSANDRA_CLIENT.execute(
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
      VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`,
      [
        ranking.year,
        ranking.makerid,
        ranking.pieid,
        ranking.userid,
        ranking.pastry,
        ranking.filling,
        ranking.topping,
        ranking.looks,
        ranking.value,
        ranking.notes,
        new Date()
      ],
      { prepare: true }
    );
    return ok(true);
  } catch {
    return err(StorageError.GenericError);
  }
}
