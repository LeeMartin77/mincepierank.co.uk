import { redirect, error } from '@sveltejs/kit';
import type { PageServerLoadEvent } from './$types';
import { getUserPieRankings, getPieByMakerAndId, type MakerPie } from '$lib/storage';

export const load = async ({ params, parent }: PageServerLoadEvent) => {
  const { year } = params;
  const { session } = await parent();
  if (!session?.user) {
    throw redirect(302, '/');
  }

  if (!session.user?.email) {
    throw error(500, 'Cannot access user info');
  }

  const yearInt = parseInt(year);

  if (!yearInt) {
    throw error(404, 'Not Found');
  }

  const userPieRankingsRes = await getUserPieRankings(yearInt, session.user.email);
  if (!userPieRankingsRes.isOk()) {
    throw error(500, 'Something went wrong');
  }

  const piesRes = await Promise.all(
    userPieRankingsRes.value.map((pr) => getPieByMakerAndId(yearInt, pr.makerid, pr.pieid))
  );

  return {
    rankings: userPieRankingsRes.value,
    pies: piesRes.reduce<MakerPie[]>((acc, res) => {
      if (res.isOk()) {
        acc.push(res.value);
      }
      return acc;
    }, [])
  };
};
