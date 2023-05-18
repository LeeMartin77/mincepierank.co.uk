import { err, ok, Result } from "neverthrow";
import { CASSANDRA_CLIENT } from "./cassandra";
import type { Maker } from "./types";
import { StorageError } from "./types";
import { rowToObject } from "./utilities";

export async function getMincePieMakers(
  client = CASSANDRA_CLIENT
): Promise<Result<Maker[], StorageError>> {
  try {
    const result = await client.execute("SELECT * FROM mincepierank.maker;");

    const mapped = result.rows.map(rowToObject);
    return ok(mapped as Maker[]);
  } catch {
    return err(StorageError.GenericError);
  }
}

export async function getMincePieMaker(
  id: string,
  client = CASSANDRA_CLIENT
): Promise<Result<Maker, StorageError>> {
  try {
    const result = await client.execute(
      "SELECT * FROM mincepierank.maker WHERE id = ?;",
      [id]
    );
    if (result.rows.length !== 1) {
      return err(StorageError.NotFound);
    }
    return ok(rowToObject(result.first()));
  } catch {
    return err(StorageError.GenericError);
  }
}
