import { getYears } from '$lib/storage/config';

export const load = async () => {
  const years = await getYears();
  return {
    years
  };
};
