import type { LayoutServerLoad } from './$types';
import { getAdmins } from '$lib/storage/config';
import { error } from '@sveltejs/kit';

export const load: LayoutServerLoad = async ({ parent }) => {
  const { session } = await parent();
  const admins = await getAdmins();
  const adminIndex = admins.findIndex((x) => x.id === session?.user?.email);

  if (!session || !session.user?.email || adminIndex === -1) {
    throw error(404, 'Not Found');
  }

  return {
    adminInfo: admins[adminIndex]
  };
};
