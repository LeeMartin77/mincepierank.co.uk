import { err, ok, Result } from "neverthrow";
import { CASSANDRA_CLIENT } from "./cassandra";
import type { MakerPieRanking } from "./types";
import { StorageError } from "./types";
import { calculateAverage, rowToObject } from "./utilities";

export async function getAllRankingsForPie(
  makerid: string,
  pieid: string,
  client = CASSANDRA_CLIENT
): Promise<Result<MakerPieRanking[], StorageError>> {
  try {
    const result = await client.execute(
      "SELECT * FROM mincepierank.maker_pie_ranking_yearly WHERE year = ? AND makerid = ? AND pieid = ? ALLOW FILTERING;",
      [2022, makerid, pieid],
      { prepare : true }
    );

    const mapped = result.rows.map(rowToObject);
    return ok(mapped as MakerPieRanking[]);
  } catch {
    return err(StorageError.GenericError);
  }
}

export async function getLatestRanking(
  client = CASSANDRA_CLIENT
): Promise<Result<MakerPieRanking, StorageError>> {
  try {
    const whenresult = await client.execute(
      "SELECT MAX(last_updated) as last_updated FROM mincepierank.maker_pie_ranking_yearly where year = 2022;"
    );

    if (whenresult.rowLength === 0) {
      return err(StorageError.NotFound);
    }

    const when = whenresult.first().get("last_updated");
    if (!when) {
      return err(StorageError.NotFound);
    }

    const top = await client.execute(
      `SELECT makerid,
      pieid,
      pastry,
      filling,
      topping,
      looks,
      value,
      CAST(last_updated as text) as last_updated FROM mincepierank.maker_pie_ranking_yearly where year = ? AND last_updated >= ? ALLOW FILTERING;`,
      [2022, new Date(when)],
      { prepare : true }
    );
    return ok(rowToObject(top.first()));
  } catch (ex: any) {
    console.error(ex.message);
    return err(StorageError.GenericError);
  }
}

export async function getMyRankingForPie(
  makerid: string,
  pieid: string,
  userid: string,
  client = CASSANDRA_CLIENT
): Promise<Result<MakerPieRanking, StorageError>> {
  try {
    const result = await client.execute(
      `SELECT makerid,
      pieid,
      pastry,
      filling,
      topping,
      looks,
      value,
      CAST(last_updated as text) as last_updated FROM mincepierank.maker_pie_ranking_yearly WHERE year = ? makerid = ? AND pieid = ? AND userid = ? ALLOW FILTERING;`,
      [2022, makerid, pieid, userid],
      { prepare : true }
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
  userid: string,
  client = CASSANDRA_CLIENT
): Promise<
  Result<
    (MakerPieRanking & { average: number; count: undefined })[],
    StorageError
  >
> {
  try {
    const result = await client.execute(
      `SELECT makerid,
        pieid,
        pastry,
        filling,
        topping,
        looks,
        value,
        CAST(last_updated as text) as last_updated
        FROM mincepierank.maker_pie_ranking_yearly WHERE year = ? AND userid = ? ALLOW FILTERING;`,
      [2022, userid],
      { prepare: true }
    );

    return ok(
      result.rows.map(rowToObject).map(addAverageScore) as (MakerPieRanking & {
        count: undefined;
        average: number;
      })[]
    );
  } catch {
    return err(StorageError.GenericError);
  }
}

export type PieRankingSummary = Omit<MakerPieRanking, "userid" | "notes"> & {
  count: number;
  average: number;
};

export function addAverageScore(mapped: any) {
  mapped.average = calculateAverage(mapped);
  return mapped;
}

export async function getPieRankingSummary(
  makerid: string,
  pieid: string,
  client = CASSANDRA_CLIENT
): Promise<Result<PieRankingSummary, StorageError>> {
  try {
    const result = await client.execute(
      `SELECT makerid,
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
      [2022, makerid, pieid],
      { prepare: true }
    );

    return ok(addAverageScore(rowToObject(result.first())));
  } catch (ex) {
    console.log(ex)
    return err(StorageError.GenericError);
  }
}

export async function getPieRankingSummariesByIds(
  makerid: string,
  pieids: string[],
  client = CASSANDRA_CLIENT
): Promise<Result<PieRankingSummary[], StorageError>> {
  try {
    const result = await client.execute(
      `SELECT makerid,
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
      [2022, makerid, pieids],
      { prepare: true }
    );

    return ok(
      result.rows.map(rowToObject).map(addAverageScore) as PieRankingSummary[]
    );
  } catch {
    return err(StorageError.GenericError);
  }
}

export async function getMakerPieRankingSummaries(
  makerid: string,
  client = CASSANDRA_CLIENT
): Promise<Result<PieRankingSummary[], StorageError>> {
  try {
    const result = await client.execute(
      `SELECT makerid,
        pieid,
        AVG(cast(pastry as float)) as pastry,
        AVG(cast(filling as float)) as filling,
        AVG(cast(topping as float)) as topping,
        AVG(cast(looks as float)) as looks,
        AVG(cast(value as float)) as value,
        CAST(MAX(last_updated) as text) as last_updated,
        CAST(COUNT(1) as int) as count 
      FROM mincepierank.maker_pie_ranking_yearly
      WHERE makerid = ?
      GROUP BY year, makerid, pieid;`,
      [makerid],
      { prepare: true }
    );

    return ok(
      result.rows.map(rowToObject).map(addAverageScore) as PieRankingSummary[]
    );
  } catch {
    return err(StorageError.GenericError);
  }
}

export async function getAllPieRankingSummaries(
  client = CASSANDRA_CLIENT
): Promise<Result<PieRankingSummary[], StorageError>> {
  try {
    const result = await client.execute(
      `SELECT makerid,
        pieid,
        AVG(cast(pastry as float)) as pastry,
        AVG(cast(filling as float)) as filling,
        AVG(cast(topping as float)) as topping,
        AVG(cast(looks as float)) as looks,
        AVG(cast(value as float)) as value,
        CAST(MAX(last_updated) as text) as last_updated,
        CAST(COUNT(1) as int) as count 
      FROM mincepierank.maker_pie_ranking_yearly
      GROUP BY year, makerid, pieid;`,
      { prepare: true }
    );

    return ok(
      result.rows.map(rowToObject).map(addAverageScore) as PieRankingSummary[]
    );
  } catch {
    return err(StorageError.GenericError);
  }
}

export async function addPieRanking(
  ranking: MakerPieRanking,
  client = CASSANDRA_CLIENT
): Promise<Result<boolean, StorageError>> {
  const mapRanking = ranking as {
    [key: string]: string | number | Date | undefined;
  };
  const rankingValues = ["pastry", "filling", "topping", "looks", "value"];

  rankingValues.forEach(
    (key) =>
      (mapRanking[key] = parseInt(
        mapRanking[key] !== undefined ? mapRanking[key]!.toString() : "NaN"
      ))
  );
  if (
    !rankingValues.every(
      (key) =>
        (mapRanking[key] as number) < 6 && (mapRanking[key] as number) > 0
    ) ||
    (ranking.notes && ranking.notes.length > 140) ||
    !ranking.makerid ||
    !ranking.pieid ||
    !ranking.userid
  ) {
    return err(StorageError.BadInput);
  }
  try {
    await client.execute(
      `INSERT INTO 
        mincepierank.maker_pie_ranking_yearly (
          year
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
        new Date(),
      ],
      { prepare: true }
    );
    return ok(true);
  } catch {
    return err(StorageError.GenericError);
  }
}
