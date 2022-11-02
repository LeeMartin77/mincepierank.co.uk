import { types } from "cassandra-driver";

export function rowToObject<T>(row: types.Row): T {
  const constructed: any = {};
  row.keys().forEach((key) => (constructed[key] = row.get(key)));
  return constructed as T;
}
