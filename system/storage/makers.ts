import { ok, Result } from "neverthrow";
import { CASSANDRA_CLIENT } from "./cassandra";
import { Maker, StorageError } from "./types";

export async function getMincePieMakers(
  client = CASSANDRA_CLIENT
): Promise<Result<Maker[], StorageError>> {
  const result = await client.execute("SELECT * FROM mincepierank.maker;");

  const mapped = result.rows.map((item) => {
    const constructed: any = {};
    item.keys().forEach((key) => (constructed[key] = item.get(key)));
    return constructed;
  });
  return ok(mapped as Maker[]);
}
