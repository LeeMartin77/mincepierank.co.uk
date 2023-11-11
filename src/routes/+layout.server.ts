import type { LayoutServerLoad } from './$types';
import { getConfig, getYears } from '$lib/storage/config';

export const load: LayoutServerLoad = async (event) => {
  const [config, years, session] = await Promise.all([
    getConfig(),
    getYears(),
    event.locals.getSession()
  ]);
  return {
    session,
    readonly: config.readonly === 'true',
    notice: config.notice,
    years,
    activeYear: parseInt(config.activeYear),
    customPiesEnabled: config.customPiesEnabled === 'true',
    imgprssr: process.env.IMGPRSSR_ROOT || 'https://static.mincepierank.co.uk'
  };
};
