import type { PageServerLoadEvent } from './$types';
import { getUserPieById, getUserPieRankingSummary } from '$lib/storage';
import { error } from '@sveltejs/kit';
export const load = async ({ params }: PageServerLoadEvent) => {
  const { year, id } = params;
  const res = await getUserPieById(parseInt(year), id);
  if (res.isErr()) {
    throw error(404, 'Not Found');
  }
  const rankingRes = await getUserPieRankingSummary(parseInt(year), id);

  return {
    pie: res.value,
    ranking: rankingRes.isOk() ? rankingRes.value : undefined
  };
};
