import { err, ok, Result } from "neverthrow";
import { CASSANDRA_CLIENT } from "./cassandra";
import { MakerPieRanking, StorageError } from "./types";
import { rowToObject } from "./utilities";

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

export async function getMyRankingForPie(
  makerid: string,
  pieid: string,
  userid: string,
  client = CASSANDRA_CLIENT
): Promise<Result<MakerPieRanking, StorageError>> {
  try {
    const result = await client.execute(
      "SELECT * FROM mincepierank.maker_pie_ranking WHERE makerid = ? AND pieid = ? AND userid = ? ALLOW FILTERING;",
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

export type PieRankingSummary = Omit<MakerPieRanking, "userid" | "notes"> & {
  count: number;
};

export async function getPieRankingSummary(
  makerid: string,
  pieid: string,
  client = CASSANDRA_CLIENT,
): Promise<Result<PieRankingSummary, StorageError>> {
  try {
    const result = await client.execute(
      `SELECT makerid,
        pieid,
        AVG(pastry) as pastry,
        AVG(filling) as filling,
        AVG(topping) as topping,
        AVG(looks) as looks,
        AVG(value) as value,
        CAST(COUNT(1) as int) as count 
      FROM mincepierank.maker_pie_ranking
      WHERE makerid = ? AND pieid = ? ALLOW FILTERING;`,
      [makerid, pieid],
      { prepare: true }
    );

    return ok(rowToObject(result.first()));
  } catch {
    return err(StorageError.GenericError);
  }
}


export async function getMakerPieRankingSummaries(
  makerid: string,
  client = CASSANDRA_CLIENT,
): Promise<Result<PieRankingSummary[], StorageError>> {
  try {
    const result = await client.execute(
      `SELECT makerid,
        pieid,
        AVG(pastry) as pastry,
        AVG(filling) as filling,
        AVG(topping) as topping,
        AVG(looks) as looks,
        AVG(value) as value,
        CAST(COUNT(1) as int) as count 
      FROM mincepierank.maker_pie_ranking
      WHERE makerid = ?;`,
      [makerid],
      { prepare: true }
    );

    return ok(result.rows.map(rowToObject) as PieRankingSummary[]);
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
        AVG(pastry) as pastry,
        AVG(filling) as filling,
        AVG(topping) as topping,
        AVG(looks) as looks,
        AVG(value) as value,
        CAST(COUNT(1) as int) as count 
      FROM mincepierank.maker_pie_ranking;`,
      { prepare: true }
    );

    return ok(result.rows.map(rowToObject) as PieRankingSummary[]);
  } catch {
    return err(StorageError.GenericError);
  }
}

export async function addPieRanking(
  ranking: MakerPieRanking,
  client = CASSANDRA_CLIENT
): Promise<Result<boolean, StorageError>> {
  const mapRanking = ranking as { [key: string]: string | number | undefined };
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
          notes
        )
      VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);`,
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
      ],
      { prepare: true }
    );
    return ok(true);
  } catch {
    return err(StorageError.GenericError);
  }
}
