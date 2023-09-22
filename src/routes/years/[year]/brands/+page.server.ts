import { error } from '@sveltejs/kit';
import type { PageServerLoadEvent } from './$types';
import { getMincePieMakersForYear } from '$lib/storage';
import { getYears } from '$lib/storage/config';

export const load = async ({ params }: PageServerLoadEvent) => {
  const years = await getYears();
  const { year } = params;
  if (!years.includes(parseInt(year))) {
    throw error(404, 'Not Found');
  }
  const yearInt = parseInt(year);
  const data = await getMincePieMakersForYear(yearInt);
  return {
    makers: data.isOk() ? data.value : [],
    year: yearInt
  };
};
