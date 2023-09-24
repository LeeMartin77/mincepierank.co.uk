import { error, json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { StorageError, addUserPieRanking, getMyRankingForUserPie } from '$lib/storage';
import { getConfig } from '$lib/storage/config';

export const GET: RequestHandler = async ({ locals, url: { searchParams } }) => {
  const session = await locals.getSession();

  const email = session?.user?.email;
  if (!email) {
    throw error(401, 'Not signed in');
  }

  const year = searchParams.get('year') || 'NaN';
  const yearInt = parseInt(year);
  const pieid = searchParams.get('pieid');

  if (!yearInt || !pieid) {
    throw error(404, 'Not Found');
  }

  const userPieRankingsRes = await getMyRankingForUserPie(yearInt, pieid, email);
  if (!userPieRankingsRes.isOk()) {
    if (userPieRankingsRes.error === StorageError.NotFound) {
      throw error(404, 'Not Ranked');
    }
    throw error(500, 'Something went wrong');
  }

  return json(userPieRankingsRes.value);
};

export const POST: RequestHandler = async ({ request, locals }) => {
  const session = await locals.getSession();
  const config = await getConfig();
  if (config.readonly === 'true') {
    throw error(400, 'Currently readonly');
  }
  const email = session?.user?.email;
  if (!email) {
    throw error(401, 'Not signed in');
  }
  const { year, pieid, pastry, filling, topping, looks, value, notes } = await request.json();
  const yearInt = parseInt(year);

  const ranking = {
    year: yearInt,
    pieid,
    userid: email,
    pastry,
    filling,
    topping,
    looks,
    value,
    notes
  };

  const res = await addUserPieRanking(ranking);
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
