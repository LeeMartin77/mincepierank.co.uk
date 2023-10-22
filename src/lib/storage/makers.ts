import { err, ok, Result } from 'neverthrow';
import { CASSANDRA_CLIENT } from './cassandra';
import type { Maker } from './types';
import { StorageError } from './types';
import { rowToObject } from './utilities';

export async function getMincePieMakers(): Promise<Result<Maker[], StorageError>> {
  try {
    const result = await CASSANDRA_CLIENT.execute('SELECT * FROM mincepierank.maker;');

    const mapped = result.rows.map(rowToObject);
    return ok(mapped as Maker[]);
  } catch {
    return err(StorageError.GenericError);
  }
}

export async function getMincePieMakersForYear(
  year: number
): Promise<Result<Maker[], StorageError>> {
  try {
    const pieResult = await CASSANDRA_CLIENT.execute(
      'SELECT makerid FROM mincepierank.maker_pie_yearly WHERE year = ?;',
      [year],
      { prepare: true }
    );
    const ids = pieResult.rows.reduce<Set<string>>((acc, row) => {
      acc.add(row.get('makerid'));
      return acc;
    }, new Set<string>());
    const result = await CASSANDRA_CLIENT.execute(
      'SELECT * FROM mincepierank.maker WHERE id in ?;',
      [Array.from(ids)],
      { prepare: true }
    );

    const mapped = result.rows.map(rowToObject);
    return ok(mapped as Maker[]);
  } catch {
    return err(StorageError.GenericError);
  }
}

export async function getMincePieMaker(id: string): Promise<Result<Maker, StorageError>> {
  try {
    const result = await CASSANDRA_CLIENT.execute(
      'SELECT * FROM mincepierank.maker WHERE id = ?;',
      [id],
      { prepare: true }
    );
    if (result.rows.length !== 1) {
      return err(StorageError.NotFound);
    }
    return ok(rowToObject(result.first()));
  } catch {
    return err(StorageError.GenericError);
  }
}
