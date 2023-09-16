import type { LayoutServerLoad } from './$types';
import { validateAdmin } from './utilities';

export const load: LayoutServerLoad = async ({ parent }) => {
  const { session } = await parent();
  const adminInfo = await validateAdmin(session?.user?.email);
  return {
    adminInfo: adminInfo
  };
};
