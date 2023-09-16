import type { PageServerLoad } from './$types';
import { getConfig, setConfigValue, deleteConfigValue } from '$lib/storage/config';
import { fail, type Actions } from '@sveltejs/kit';
import { validateAdmin } from '../utilities';

export const load: PageServerLoad = async () => {
  const config = await getConfig();
  return {
    config
  };
};

export const actions = {
  set: async ({ request, locals }) => {
    const session = await locals.getSession();
    await validateAdmin(session?.user?.email);
    const data = await request.formData();
    const key = data.get('key');
    const value = data.get('value');
    if (!key || !value) {
      return fail(400, { key, value, incorrect: true });
    }

    await setConfigValue(key as string, value as string);

    return { success: true };
  },
  delete: async ({ request, locals }) => {
    const session = await locals.getSession();
    await validateAdmin(session?.user?.email);
    const data = await request.formData();
    const key = data.get('key');
    if (!key) {
      return fail(400, { key, incorrect: true });
    }

    await deleteConfigValue(key as string);

    return { success: true };
  }
} satisfies Actions;
