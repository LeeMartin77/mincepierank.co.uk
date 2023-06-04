import { redirect, error } from '@sveltejs/kit';
import type { PageServerLoadEvent } from './$types';
import { getUserPieRankings, getAllMakerPies } from '$lib/storage';

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
  const piesRes = await getAllMakerPies(yearInt);
  if (!userPieRankingsRes.isOk() || !piesRes.isOk()) {
    throw error(500, 'Something went wrong');
  }
  return {
    rankings: userPieRankingsRes.value,
    pies: piesRes.value
  };
};
