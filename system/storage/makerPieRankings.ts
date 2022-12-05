import { err, ok, Result } from "neverthrow";
import { CASSANDRA_CLIENT } from "./cassandra";
import { MakerPieRanking, StorageError } from "./types";
import { calculateAverage, rowToObject } from "./utilities";

export async function getAllRankingsForPie(
  makerid: string,
  pieid: string,
  client = CASSANDRA_CLIENT
): Promise<Result<MakerPieRanking[], StorageError>> {
  try {
    const result = await client.execute(
      "SELECT * FROM mincepierank.maker_pie_ranking WHERE makerid = ? AND pieid = ? ALLOW FILTERING;",
      [makerid, pieid]
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
      "SELECT MAX(last_updated) as last_updated FROM mincepierank.maker_pie_ranking;"
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
      CAST(last_updated as text) as last_updated FROM mincepierank.maker_pie_ranking where last_updated >= ? ALLOW FILTERING;`,
      [new Date(when)]
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
      CAST(last_updated as text) as last_updated FROM mincepierank.maker_pie_ranking WHERE makerid = ? AND pieid = ? AND userid = ? ALLOW FILTERING;`,
      [makerid, pieid, userid]
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
        FROM mincepierank.maker_pie_ranking WHERE userid = ? ALLOW FILTERING;`,
      [userid],
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
      FROM mincepierank.maker_pie_ranking
      WHERE makerid = ? AND pieid = ? ALLOW FILTERING;`,
      [makerid, pieid],
      { prepare: true }
    );

    return ok(addAverageScore(rowToObject(result.first())));
  } catch {
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
      FROM mincepierank.maker_pie_ranking
      WHERE makerid = ? and pieid in ?
      GROUP BY makerid, pieid;`,
      [makerid, pieids],
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
      FROM mincepierank.maker_pie_ranking
      WHERE makerid = ?
      GROUP BY makerid, pieid;`,
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
      FROM mincepierank.maker_pie_ranking
      GROUP BY makerid, pieid;`,
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
        mincepierank.maker_pie_ranking (
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
      VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`,
      [
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
