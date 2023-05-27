import { error } from '@sveltejs/kit';
import type { PageServerLoadEvent } from './$types';
import { getMakerPieRankingSummaries, getMincePieMaker, getPiesByMaker } from '$lib/storage';

export const load = async ({ params }: PageServerLoadEvent) => {
  const { year, brand } = params;
  const makerData = await getMincePieMaker(brand);
  const pies = await getPiesByMaker(parseInt(year), brand);
  if (makerData.isErr() || pies.isErr()) {
    throw error(500, 'Something went wrong');
  }
  if (pies.value.length === 0) {
    throw error(404, 'Maker not found for year');
  }

  const categories = Array.from(
    pies.value.reduce<Set<string>>((acc, curr) => {
      curr.labels.forEach((label) => {
        acc.add(label);
      });
      return acc;
    }, new Set<string>())
  );

  const pieRankings = await getMakerPieRankingSummaries(parseInt(year), brand);
  return {
    makerData: makerData.value,
    pies: pies.value,
    pieRankings: pieRankings.isOk() ? pieRankings.value : [],
    year: parseInt(year),
    categories
  };
};
