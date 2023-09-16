import { CASSANDRA_CLIENT } from './cassandra';
import { rowToObject } from './utilities';

export async function getConfig(): Promise<{ [key: string]: string }> {
  const result = await CASSANDRA_CLIENT.execute('SELECT key, value FROM mincepierank.config;');

  const mapped = result.rows.map((x) => rowToObject<{ key: string; value: string }>(x));
  const config: { [key: string]: string } = {};
  for (const val of mapped) {
    config[val.key] = val.value;
  }
  return config;
}

export async function getYears(): Promise<number[]> {
  const result = await CASSANDRA_CLIENT.execute('SELECT year FROM mincepierank.year;');
  return result.rows.map((x) => parseInt(x.get('year')));
}
