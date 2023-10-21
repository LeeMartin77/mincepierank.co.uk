import type { PageServerLoadEvent } from './$types';
import { getMincePieMaker, getPieByMakerAndId, getPieRankingSummary } from '$lib/storage';
import { error } from '@sveltejs/kit';
export const load = async ({ params }: PageServerLoadEvent) => {
  const { year, brand, pie } = params;
  const res = await getPieByMakerAndId(parseInt(year), brand, pie);
  if (res.isErr()) {
    throw error(404, 'Not Found');
  }
  const maker = await getMincePieMaker(res.value.makerid);
  const rankingRes = await getPieRankingSummary(parseInt(year), brand, pie);

  return {
    pie: res.value,
    maker: maker.isErr() ? undefined : maker.value,
    ranking: rankingRes.isOk() ? rankingRes.value : undefined
  };
};
