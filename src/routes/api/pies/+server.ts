import { error, json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { StorageError, getPiesByMaker } from '$lib/storage';

export const GET: RequestHandler = async ({ url: { searchParams } }) => {
  const year = searchParams.get('year') || 'NaN';
  const yearInt = parseInt(year);
  const makerid = searchParams.get('makerid');

  if (!yearInt || !makerid) {
    throw error(404, 'Not Found');
  }

  const userPieRankingsRes = await getPiesByMaker(yearInt, makerid);
  if (!userPieRankingsRes.isOk()) {
    if (userPieRankingsRes.error === StorageError.NotFound) {
      throw error(404, 'Not Ranked');
    }
    throw error(500, 'Something went wrong');
  }

  return json(userPieRankingsRes.value);
};
