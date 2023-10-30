import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { getAnonId } from '$lib/anonid';
import { claimAnonymous } from '$lib/claimAnonymous';

export const POST: RequestHandler = async ({ cookies, locals }) => {
  const session = await locals.getSession();
  const userid = session?.user?.email;
  const anonUserId = getAnonId(cookies);
  if (!userid) {
    // not signed in, nothing to do
    return json({});
  }

  return json(await claimAnonymous(userid, anonUserId));
};
