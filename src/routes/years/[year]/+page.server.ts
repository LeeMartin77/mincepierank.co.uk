import { error } from '@sveltejs/kit';
import type { PageServerLoadEvent } from './$types';
import { getAllMakerPies, getMincePieMakersForYear, type MakerPie } from '$lib/storage';

// TODO: Make this truly data driven
const years = ['2022', '2023'];

export const load = async ({ params }: PageServerLoadEvent) => {
  const { year } = params;
  if (!years.includes(year)) {
    throw error(404, 'Not Found');
  }
  const yearInt = parseInt(year);
  const data = await getMincePieMakersForYear(yearInt);
  const piesRes = await getAllMakerPies(yearInt);
  let pies: MakerPie[] = [];
  if (piesRes.isOk()) {
    pies = piesRes.value;
  }
  const categories = Array.from(
    pies.reduce<Set<string>>((acc, curr) => {
      curr.labels.forEach((label) => {
        acc.add(label);
      });
      return acc;
    }, new Set<string>())
  );
  return {
    makers: data.isOk() ? data.value : [],
    categories,
    year: yearInt
  };
};
