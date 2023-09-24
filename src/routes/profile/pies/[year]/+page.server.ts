import { redirect, error } from '@sveltejs/kit';
import type { PageServerLoadEvent } from './$types';
import {
  getUserPieRankings,
  getPieByMakerAndId,
  type MakerPie,
  getUserPieUserRankings,
  getUserPieById,
  type UserPie,
  getUserPiesByUser
} from '$lib/storage';

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

  const customPieRankings = await getUserPieUserRankings(yearInt, session.user.email);
  const customPies = await getUserPiesByUser(yearInt, session.user.email);
  if (!customPies.isOk() || !customPieRankings.isOk()) {
    throw error(500, 'Something went wrong');
  }
  return {
    customRankings: customPieRankings.value,
    customPies: customPies.value
  };
};
