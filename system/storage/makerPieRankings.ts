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
