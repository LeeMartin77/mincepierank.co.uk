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

export type PieRankingSummary = Omit<MakerPieRanking, "userid" | "notes"> & {
  count: number;
};

export async function getPieRankingSummary(
  makerid: string,
  pieid: string,
  client = CASSANDRA_CLIENT,
  rowToObjectFn = rowToObject
): Promise<Result<PieRankingSummary, StorageError>> {
  try {
    const result = await client.execute(
      "SELECT * FROM mincepierank.maker_pie_ranking WHERE makerid = ? AND pieid = ? ALLOW FILTERING;",
      [makerid, pieid]
    );

    const mapped = result.rows.map(rowToObjectFn) as MakerPieRanking[];
    const reduced: PieRankingSummary = mapped.reduce(
      (prev, curr) => {
        return {
          ...prev,
          pastry: curr.pastry + prev.pastry,
          filling: curr.filling + prev.filling,
          topping: curr.topping + prev.topping,
          looks: curr.looks + prev.looks,
          value: curr.value + prev.value,
          count: prev.count + 1,
        };
      },
      {
        makerid,
        pieid,
        pastry: 0,
        filling: 0,
        topping: 0,
        looks: 0,
        value: 0,
        count: 0,
      }
    );
    if (reduced.count > 1) {
      reduced.pastry = reduced.pastry / reduced.count;
      reduced.filling = reduced.filling / reduced.count;
      reduced.topping = reduced.topping / reduced.count;
      reduced.looks = reduced.looks / reduced.count;
      reduced.value = reduced.value / reduced.count;
    }
    return ok(reduced);
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
      ]
    );
    return ok(true);
  } catch {
    return err(StorageError.GenericError);
  }
}
