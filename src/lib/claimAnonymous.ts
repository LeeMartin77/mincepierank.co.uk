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
import { getYears } from '$lib/storage/config';

export async function claimAnonymous(userid: string, anonUserId: string) {
  const years = await getYears();
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
    return {};
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

  return {};
}
