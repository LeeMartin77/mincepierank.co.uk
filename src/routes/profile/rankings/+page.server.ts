import { redirect } from '@sveltejs/kit';
import type { PageServerLoadEvent } from './$types';

export const load = async ({ parent }: PageServerLoadEvent) => {
  const { session } = await parent();
  if (!session?.user) {
    throw redirect(302, '/');
  }
  return {};
};
