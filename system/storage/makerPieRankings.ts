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
