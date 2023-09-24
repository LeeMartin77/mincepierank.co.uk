import { err, ok, Result } from 'neverthrow';
import { CASSANDRA_CLIENT } from './cassandra';
import type { UserPie } from './types';
import { StorageError } from './types';
import { rowToObject } from './utilities';

const SAFE_FIELDS = `year,
id,
maker,
location,
displayname,
fresh,
labels,
image_file,
web_link,
pack_count,
pack_price_in_pence,
clean`;

export async function getUserPiesByUser(
  year: number,
  userid: string
): Promise<Result<UserPie[], StorageError>> {
  try {
    const result = await CASSANDRA_CLIENT.execute(
      `SELECT ${SAFE_FIELDS} FROM mincepierank.user_pie_yearly WHERE year = ? AND owner_userid = ? ALLOW FILTERING`,
      [year, userid],
      { prepare: true }
    );

    const mapped = result.rows.map(rowToObject) as UserPie[];
    return ok(
      mapped.map((x) => {
        delete x.owner_userid;
        return x;
      })
    );
  } catch {
    return err(StorageError.GenericError);
  }
}

export async function getUserPieById(
  year: number,
  id: string
): Promise<Result<UserPie, StorageError>> {
  try {
    const result = await CASSANDRA_CLIENT.execute(
      `SELECT ${SAFE_FIELDS} FROM mincepierank.user_pie_yearly WHERE year = ? AND id = ?;`,
      [year, id],
      { prepare: true }
    );
    if (result.rows.length !== 1) {
      return err(StorageError.NotFound);
    }
    return ok(rowToObject(result.first()) as UserPie);
  } catch {
    return err(StorageError.GenericError);
  }
}

export async function setUserPie(pie: UserPie): Promise<Result<boolean, StorageError>> {
  if (!pie.owner_userid) {
    return err(StorageError.BadInput);
  }
  try {
    await CASSANDRA_CLIENT.execute(
      `INSERT INTO mincepierank.user_pie_yearly
        (
          year,
          owner_userid,
          id,
          maker,
          location,
          displayname,
          fresh,
          labels,
          image_file,
          web_link,
          pack_count,
          pack_price_in_pence,
          clean
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
          ?,
          ?,
          ?,
          ?
        )
      `,
      [
        pie.year,
        pie.owner_userid,
        pie.id,
        pie.maker,
        pie.location,
        pie.displayname,
        pie.fresh,
        pie.labels,
        pie.image_file,
        pie.web_link,
        pie.pack_count,
        pie.pack_price_in_pence,
        pie.clean
      ],
      { prepare: true }
    );
    return ok(true);
  } catch (ex) {
    console.error(ex);
    return err(StorageError.GenericError);
  }
}
