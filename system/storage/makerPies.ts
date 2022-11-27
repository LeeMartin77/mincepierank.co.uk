import { err, ok, Result } from "neverthrow";
import { CASSANDRA_CLIENT } from "./cassandra";
import { MakerPie, StorageError } from "./types";
import { rowToObject } from "./utilities";

export async function getAllMakerPies(
  client = CASSANDRA_CLIENT
): Promise<Result<MakerPie[], StorageError>> {
  try {
    const result = await client.execute(
      "SELECT * FROM mincepierank.maker_pie;"
    );

    const mapped = result.rows.map(rowToObject);
    return ok(mapped as MakerPie[]);
  } catch {
    return err(StorageError.GenericError);
  }
}

export async function getPiesByMaker(
  makerId: string,
  client = CASSANDRA_CLIENT
): Promise<Result<MakerPie[], StorageError>> {
  try {
    const result = await client.execute(
      "SELECT * FROM mincepierank.maker_pie WHERE makerId = ?;",
      [makerId]
    );

    const mapped = result.rows.map(rowToObject);
    return ok(mapped as MakerPie[]);
  } catch {
    return err(StorageError.GenericError);
  }
}

export async function getPiesWithCategory(
  category: string,
  client = CASSANDRA_CLIENT
): Promise<Result<MakerPie[], StorageError>> {
  try {
    const result = await client.execute(
      "SELECT * FROM mincepierank.maker_pie WHERE labels contains ?;",
      [category],
      { prepare: true }
    );

    const mapped = result.rows.map(rowToObject);
    return ok(mapped as MakerPie[]);
  } catch {
    return err(StorageError.GenericError);
  }
}

export async function getPieByMakerAndId(
  makerId: string,
  id: string,
  client = CASSANDRA_CLIENT
): Promise<Result<MakerPie, StorageError>> {
  try {
    const result = await client.execute(
      "SELECT * FROM mincepierank.maker_pie WHERE makerId = ? AND id = ?;",
      [makerId, id]
    );
    if (result.rows.length !== 1) {
      return err(StorageError.NotFound);
    }
    return ok(rowToObject(result.first()) as MakerPie);
  } catch {
    return err(StorageError.GenericError);
  }
}
