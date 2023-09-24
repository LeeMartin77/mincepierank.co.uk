import { error, json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { StorageError, addPieRanking, getMyRankingForPie } from '$lib/storage';
import { getAnonId } from '$lib/anonid';
import { getConfig } from '$lib/storage/config';

export const GET: RequestHandler = async ({ cookies, locals, url: { searchParams } }) => {
  const session = await locals.getSession();

  let userid = session?.user?.email;
  if (!userid) {
    userid = getAnonId(cookies);
  }

  const year = searchParams.get('year') || 'NaN';
  const yearInt = parseInt(year);
  const makerid = searchParams.get('makerid');
  const pieid = searchParams.get('pieid');

  if (!yearInt || !makerid || !pieid) {
    throw error(404, 'Not Found');
  }

  const userPieRankingsRes = await getMyRankingForPie(yearInt, makerid, pieid, userid);
  if (!userPieRankingsRes.isOk()) {
    if (userPieRankingsRes.error === StorageError.NotFound) {
      throw error(404, 'Not Ranked');
    }
    throw error(500, 'Something went wrong');
  }

  return json(userPieRankingsRes.value);
};

export const POST: RequestHandler = async ({ cookies, request, locals }) => {
  const session = await locals.getSession();
  const config = await getConfig();
  if (config.readonly === 'true') {
    throw error(400, 'Currently readonly');
  }
  let userid = session?.user?.email;
  if (!userid) {
    userid = getAnonId(cookies);
  }

  const { year, makerid, pieid, pastry, filling, topping, looks, value, notes } =
    await request.json();
  const yearInt = parseInt(year);

  const ranking = {
    year: yearInt,
    makerid,
    pieid,
    userid,
    pastry,
    filling,
    topping,
    looks,
    value,
    notes
  };

  const res = await addPieRanking(ranking);
  if (!res.isOk()) {
    if (res.error === StorageError.NotFound) {
      throw error(404, 'Not Ranked');
    }
    if (res.error === StorageError.BadInput) {
      throw error(403, 'Bad Input');
    }
    throw error(500, 'Something went wrong');
  }

  return json(ranking);
};
