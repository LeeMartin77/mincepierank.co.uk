import { error } from '@sveltejs/kit';
import type { PageServerLoadEvent } from './$types';
import { getMincePieMaker, getPiesByMaker } from '$lib/storage';

export const load = (async ({ params }: PageServerLoadEvent) => {
  const { year, maker } = params;
  const makerData = (await getMincePieMaker(maker)).unwrapOr(undefined);
  const pies = (await getPiesByMaker(parseInt(year), maker)).unwrapOr([]);
  if (makerData === undefined || pies.length === 0) {
    throw error(404, 'Maker not found for year');
  }
  return {
    makerData,
    pies
  };
});