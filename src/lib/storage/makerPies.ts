import { err, ok, Result } from 'neverthrow';
import { CASSANDRA_CLIENT } from './cassandra';
import type { MakerPie } from './types';
import { StorageError } from './types';
import { rowToObject } from './utilities';

export async function getAllMakerPies(year: number): Promise<Result<MakerPie[], StorageError>> {
  try {
    const result = await CASSANDRA_CLIENT.execute(
      'SELECT * FROM mincepierank.maker_pie_yearly where year = ?;',
      [year],
      { prepare: true }
    );

    const mapped = result.rows.map(rowToObject);
    return ok(mapped as MakerPie[]);
  } catch {
    return err(StorageError.GenericError);
  }
}

export async function getPiesByMaker(
  year: number,
  makerId: string
): Promise<Result<MakerPie[], StorageError>> {
  try {
    const result = await CASSANDRA_CLIENT.execute(
      'SELECT * FROM mincepierank.maker_pie_yearly WHERE year = ? AND makerid = ? ALLOW FILTERING',
      [year, makerId],
      { prepare: true }
    );

    const mapped = result.rows.map(rowToObject);
    return ok(mapped as MakerPie[]);
  } catch {
    return err(StorageError.GenericError);
  }
}

export async function getPiesWithCategory(
  year: number,
  category: string
): Promise<Result<MakerPie[], StorageError>> {
  try {
    const result = await CASSANDRA_CLIENT.execute(
      'SELECT * FROM mincepierank.maker_pie_yearly WHERE year = ? AND labels contains ? ALLOW FILTERING;',
      [year, category],
      { prepare: true }
    );

    const mapped = result.rows.map(rowToObject);
    return ok(mapped as MakerPie[]);
  } catch {
    return err(StorageError.GenericError);
  }
}

export async function getPieByMakerAndId(
  year: number,
  makerId: string,
  id: string
): Promise<Result<MakerPie, StorageError>> {
  try {
    const result = await CASSANDRA_CLIENT.execute(
      'SELECT * FROM mincepierank.maker_pie_yearly WHERE year = ? AND makerId = ? AND id = ?;',
      [year, makerId, id],
      { prepare: true }
    );
    if (result.rows.length !== 1) {
      return err(StorageError.NotFound);
    }
    return ok(rowToObject(result.first()) as MakerPie);
  } catch {
    return err(StorageError.GenericError);
  }
}
