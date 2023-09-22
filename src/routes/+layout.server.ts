import type { LayoutServerLoad } from './$types';
import { getConfig, getYears } from '$lib/storage/config';

export const load: LayoutServerLoad = async (event) => {
  const config = await getConfig();
  return {
    session: await event.locals.getSession(),
    readonly: config.readonly === 'true',
    years: await getYears(),
    activeYear: parseInt(config.activeYear),
    imgprssr: process.env.IMGPRSSR_ROOT || 'https://static.mincepierank.co.uk'
  };
};
