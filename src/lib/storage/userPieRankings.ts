import { err, ok, Result } from 'neverthrow';
import { CASSANDRA_CLIENT } from './cassandra';
import type { UserPieRanking } from './types';
import { StorageError } from './types';
import { calculateAverage, rowToObject } from './utilities';

const SAFE_FIELDS = `year,
pieid,
pastry,
filling,
topping,
looks,
value,
CAST(last_updated as text)`;

export async function getAllRankingsForUserPie(
  year: number,
  pieid: string
): Promise<Result<UserPieRanking[], StorageError>> {
  try {
    const result = await CASSANDRA_CLIENT.execute(
      `SELECT ${SAFE_FIELDS} FROM mincepierank.user_pie_ranking_yearly WHERE year = ? AND pieid = ? ALLOW FILTERING;`,
      [year, pieid],
      { prepare: true }
    );

    const mapped = result.rows.map(rowToObject);
    return ok(mapped as UserPieRanking[]);
  } catch {
    return err(StorageError.GenericError);
  }
}

export async function getMyRankingForUserPie(
  year: number,
  pieid: string,
  userid: string
): Promise<Result<UserPieRanking, StorageError>> {
  try {
    const result = await CASSANDRA_CLIENT.execute(
      `SELECT ${SAFE_FIELDS} as last_updated FROM mincepierank.user_pie_ranking_yearly WHERE year = ? AND pieid = ? AND userid = ? ALLOW FILTERING;`,
      [year, pieid, userid],
      { prepare: true }
    );
    if (result.rows.length === 0) {
      return err(StorageError.NotFound);
    }

    return ok(rowToObject(result.first()) as UserPieRanking);
  } catch (ex) {
    console.error(ex);
    return err(StorageError.GenericError);
  }
}

export async function getUserPieUserRankings(
  year: number,
  userid: string
): Promise<Result<(UserPieRanking & { average: number; count: undefined })[], StorageError>> {
  try {
    const result = await CASSANDRA_CLIENT.execute(
      `SELECT year,
        pieid,
        pastry,
        filling,
        topping,
        looks,
        value,
        CAST(last_updated as text) as last_updated
        FROM mincepierank.user_pie_ranking_yearly WHERE year = ? AND userid = ? ALLOW FILTERING;`,
      [year, userid],
      { prepare: true }
    );

    return ok(
      result.rows.map<UserPieRanking>(rowToObject).map(addAverageScore) as (UserPieRanking & {
        count: undefined;
        average: number;
      })[]
    );
  } catch {
    return err(StorageError.GenericError);
  }
}

export type UserPieRankingSummary = Omit<UserPieRanking, 'userid' | 'notes'> & {
  count: number | undefined;
  average: number;
};

function addAverageScore<T>(
  mapped: Pick<UserPieRanking, 'filling' | 'pastry' | 'topping' | 'value' | 'looks'> & T
) {
  return { ...mapped, average: calculateAverage(mapped) };
}

export async function getUserPieRankingSummary(
  year: number,
  pieid: string
): Promise<Result<UserPieRankingSummary, StorageError>> {
  try {
    const result = await CASSANDRA_CLIENT.execute(
      `SELECT year,
        pieid,
        AVG(cast(pastry as float)) as pastry,
        AVG(cast(filling as float)) as filling,
        AVG(cast(topping as float)) as topping,
        AVG(cast(looks as float)) as looks,
        AVG(cast(value as float)) as value,
        CAST(MAX(last_updated) as text) as last_updated,
        CAST(COUNT(1) as int) as count 
      FROM mincepierank.user_pie_ranking_yearly
      WHERE year = ? AND pieid = ? ALLOW FILTERING;`,
      [year, pieid],
      { prepare: true }
    );

    return ok(addAverageScore(rowToObject(result.first())));
  } catch (ex) {
    console.log(ex);
    return err(StorageError.GenericError);
  }
}

export async function getUserPieRankingSummariesByIds(
  year: number,
  pieids: string[]
): Promise<Result<UserPieRankingSummary[], StorageError>> {
  try {
    const result = await CASSANDRA_CLIENT.execute(
      `SELECT year,
        pieid,
        AVG(cast(pastry as float)) as pastry,
        AVG(cast(filling as float)) as filling,
        AVG(cast(topping as float)) as topping,
        AVG(cast(looks as float)) as looks,
        AVG(cast(value as float)) as value,
        CAST(MAX(last_updated) as text) as last_updated,
        CAST(COUNT(1) as int) as count 
      FROM mincepierank.user_pie_ranking_yearly
      WHERE year = ? and pieid in ?
      GROUP BY pieid;`,
      [year, pieids],
      { prepare: true }
    );

    return ok(result.rows.map<UserPieRankingSummary>(rowToObject).map(addAverageScore));
  } catch {
    return err(StorageError.GenericError);
  }
}

export async function addUserPieRanking(
  ranking: UserPieRanking
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
    !ranking.pieid ||
    !ranking.userid ||
    !ranking.year
  ) {
    return err(StorageError.BadInput);
  }
  try {
    await CASSANDRA_CLIENT.execute(
      `INSERT INTO 
        mincepierank.user_pie_ranking_yearly (
          year,
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
      VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`,
      [
        ranking.year,
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
