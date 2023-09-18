import { err, ok, Result } from 'neverthrow';
import { CASSANDRA_CLIENT } from './cassandra';
import type { MakerPie } from './types';
import { StorageError } from './types';
import { rowToObject } from './utilities';

export async function getAllMakerPies(
  year: number,
  validated = true
): Promise<Result<MakerPie[], StorageError>> {
  try {
    const result = await CASSANDRA_CLIENT.execute(
      'SELECT * FROM mincepierank.maker_pie_yearly where year = ? AND validated = ? ALLOW FILTERING',
      [year, validated],
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
  makerId: string,
  validated = true
): Promise<Result<MakerPie[], StorageError>> {
  try {
    const result = await CASSANDRA_CLIENT.execute(
      'SELECT * FROM mincepierank.maker_pie_yearly WHERE year = ? AND makerid = ? AND validated = ? ALLOW FILTERING',
      [year, makerId, validated],
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
  category: string,
  validated = true
): Promise<Result<MakerPie[], StorageError>> {
  try {
    const result = await CASSANDRA_CLIENT.execute(
      'SELECT * FROM mincepierank.maker_pie_yearly WHERE year = ? AND validated = ? AND labels contains ? ALLOW FILTERING;',
      [year, category, validated],
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

export async function setMakerPie(pie: MakerPie): Promise<Result<boolean, StorageError>> {
  try {
    await CASSANDRA_CLIENT.execute(
      `INSERT INTO mincepierank.maker_pie_yearly
        (
          'year',
          'makerid',
          'id',
          'displayname',
          'fresh',
          'labels',
          'image_file',
          'web_link',
          'pack_count', 
          'pack_price_in_pence'
        )
        VALUES
        (
          ?,
          ?,
          ?,
          ?,
          ?,
          ?,
          ?,
          ?,
          ?,
          ?
        )
      `,
      [
        pie.year,
        pie.makerid,
        pie.id,
        pie.displayname,
        pie.fresh,
        pie.labels,
        pie.image_file,
        pie.web_link,
        pie.pack_count,
        pie.pack_price_in_pence
      ],
      { prepare: true }
    );
    return ok(true);
  } catch (ex) {
    console.error(ex);
    return err(StorageError.GenericError);
  }
}
