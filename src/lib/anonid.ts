import type { Cookies } from '@sveltejs/kit';
import { randomUUID } from 'node:crypto';

const isGuid = (str: string) => {
  const guidPattern = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i;
  return guidPattern.test(str);
};

export const getAnonId = (cookies: Cookies): string => {
  const anonid = cookies.get('anonid');
  if (!anonid || !isGuid(anonid)) {
    const newId = randomUUID();
    cookies.set('anonid', newId, { path: '/' });
    return newId;
  }
  return anonid;
};
