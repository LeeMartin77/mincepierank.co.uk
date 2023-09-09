import type { LayoutServerLoad } from './$types';
import { env } from '$env/dynamic/private';

export const load: LayoutServerLoad = async (event) => {
  return {
    session: await event.locals.getSession(),
    readonly: env.READONLY === 'true'
  };
};
