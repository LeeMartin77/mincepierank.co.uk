import { error } from '@sveltejs/kit';
import type { PageServerLoadEvent } from './$types';
import { getAllMakerPies, getAllPieRankingSummaries } from '$lib/storage';

export const load = async ({ params }: PageServerLoadEvent) => {
  const { year } = params;
  const pies = await getAllMakerPies(parseInt(year));
  if (pies.isErr()) {
    throw error(500, 'Something went wrong');
  }
  if (pies.value.length === 0) {
    throw error(404, 'No pies for this year');
  }

  const categories = Array.from(
    pies.value.reduce<Set<string>>((acc, curr) => {
      curr.labels.forEach((label) => {
        acc.add(label);
      });
      return acc;
    }, new Set<string>())
  );

  const pieRankings = await getAllPieRankingSummaries(parseInt(year));
  return {
    pies: pies.value,
    pieRankings: pieRankings.isOk() ? pieRankings.value : [],
    year: parseInt(year),
    categories
  };
};
