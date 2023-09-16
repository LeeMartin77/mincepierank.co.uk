import type { LayoutServerLoad } from './$types';
import { getConfig } from '$lib/storage/config';

export const load: LayoutServerLoad = async (event) => {
  const config = await getConfig();
  return {
    session: await event.locals.getSession(),
    readonly: config.readonly === 'true'
  };
};
