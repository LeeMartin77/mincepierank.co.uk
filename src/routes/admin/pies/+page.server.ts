import type { PageServerLoad } from './$types';
import { getConfig, getYears } from '$lib/storage/config';

import { getAllMakerPies, type MakerPie } from '$lib/storage';

export const load: PageServerLoad = async (event) => {
  const config = await getConfig();
  const year = parseInt(event.url.searchParams.get('year') || '0') || new Date().getFullYear();
  const years = await getYears();
  const validatedPies = (await getAllMakerPies(year, true)).unwrapOr([] as MakerPie[]);
  const unvalidatedPies = (await getAllMakerPies(year, false)).unwrapOr([] as MakerPie[]);
  return {
    config,
    year,
    years,
    validatedPies,
    unvalidatedPies
  };
};
