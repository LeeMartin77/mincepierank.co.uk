import { getAdmins } from '$lib/storage/config';
import { error } from '@sveltejs/kit';

export const validateAdmin = async (id: string | null | undefined) => {
  if (!id) {
    throw error(404, 'Not Found');
  }
  const admins = await getAdmins();
  const adminIndex = admins.findIndex((x) => x.id === id);

  if (adminIndex === -1) {
    throw error(404, 'Not Found');
  }
  return admins[adminIndex];
};
