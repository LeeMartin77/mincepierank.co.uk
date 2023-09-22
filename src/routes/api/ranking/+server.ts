import { error, json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { StorageError, addPieRanking, getMyRankingForPie } from '$lib/storage';

export const GET: RequestHandler = async ({ locals, url: { searchParams } }) => {
  const session = await locals.getSession();

  const email = session?.user?.email;
  if (!email) {
    throw error(401, 'Not signed in');
  }

  const year = searchParams.get('year') || 'NaN';
  const yearInt = parseInt(year);
  const makerid = searchParams.get('makerid');
  const pieid = searchParams.get('pieid');

  if (!yearInt || !makerid || !pieid) {
    throw error(404, 'Not Found');
  }

  const userPieRankingsRes = await getMyRankingForPie(yearInt, makerid, pieid, email);
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
  const email = session?.user?.email;
  if (!email) {
    throw error(401, 'Not signed in');
  }
  const { year, makerid, pieid, pastry, filling, topping, looks, value, notes } =
    await request.json();
  const yearInt = parseInt(year);

  await addPieRanking({
    year,
    makerid,
    pieid,
    userid: email,
    pastry,
    filling,
    topping,
    looks,
    value,
    notes
  });

  const userPieRankingsRes = await getMyRankingForPie(yearInt, makerid, pieid, email);
  if (!userPieRankingsRes.isOk()) {
    if (userPieRankingsRes.error === StorageError.NotFound) {
      throw error(404, 'Not Ranked');
    }
    if (userPieRankingsRes.error === StorageError.BadInput) {
      throw error(403, 'Bad Input');
    }
    throw error(500, 'Something went wrong');
  }

  return json(userPieRankingsRes.value);
};
