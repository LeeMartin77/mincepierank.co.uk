import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import {
  addPieRanking,
  getUserPieRankings,
  getUserPieUserRankings,
  type UserPieRanking,
  type MakerPieRanking,
  deleteRanking,
  addUserPieRanking,
  deleteUserRanking
} from '$lib/storage';
import { getAnonId } from '$lib/anonid';
import { getYears } from '$lib/storage/config';

export const POST: RequestHandler = async ({ cookies, locals }) => {
  const session = await locals.getSession();
  const years = await getYears();
  const userid = session?.user?.email;
  const anonUserId = getAnonId(cookies);
  if (!userid) {
    // not signed in, nothing to do
    return json({});
  }

  const allRankings = (
    await Promise.all(
      years.map((x) =>
        getUserPieRankings(x, anonUserId).then((res) => {
          if (res.isOk()) {
            return res.value;
          }
          return [];
        })
      )
    )
  ).reduce<
    (MakerPieRanking & {
      average: number;
      count: undefined;
    })[]
  >((acc, curr) => {
    return [...acc, ...curr];
  }, []);

  const allUserRankings = (
    await Promise.all(
      years.map((x) =>
        getUserPieUserRankings(x, anonUserId).then((res) => {
          if (res.isOk()) {
            return res.value;
          }
          return [];
        })
      )
    )
  ).reduce<
    (UserPieRanking & {
      average: number;
      count: undefined;
    })[]
  >((acc, curr) => {
    return [...acc, ...curr];
  }, []);

  if (allRankings.length === 0 && allUserRankings.length === 0) {
    // No anon pies to move
    return json({});
  }

  await Promise.all(
    allRankings.map(async (x) => {
      x.userid = userid;
      const create = await addPieRanking(x);
      if (create.isOk()) {
        await deleteRanking({ ...x, userid: anonUserId });
      }
    })
  );

  await Promise.all(
    allUserRankings.map(async (x) => {
      x.userid = userid;
      const create = await addUserPieRanking(x);
      if (create.isOk()) {
        await deleteUserRanking({ ...x, userid: anonUserId });
      }
    })
  );

  return json({});
};
