import { error } from '@sveltejs/kit';
import type { PageServerLoadEvent } from './$types';
import { getMincePieMakersForYear } from '$lib/storage';
import { getConfig } from '$lib/storage/config';

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export const load = async (_event: PageServerLoadEvent) => {
  const config = await getConfig();
  if (config.readonly === 'true') {
    throw error(404, 'Not Found');
  }
  const activeYear = parseInt(config.activeYear);
  const data = await getMincePieMakersForYear(activeYear);
  return {
    makers: data.isOk() ? data.value : [],
    activeYear
  };
};
