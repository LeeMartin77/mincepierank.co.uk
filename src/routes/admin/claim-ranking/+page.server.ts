import { fail, type Actions } from '@sveltejs/kit';
import { validateAdmin } from '../utilities';
import { claimAnonymous } from '$lib/claimAnonymous';

export const actions = {
  claim: async ({ request, locals }) => {
    const session = await locals.getSession();
    await validateAdmin(session?.user?.email);
    const data = await request.formData();
    const userid = data.get('userid');
    const anonuserid = data.get('anonuserid');
    if (!userid || !anonuserid) {
      return fail(400, { userid, anonuserid, incorrect: true });
    }

    return await claimAnonymous(userid as string, anonuserid as string);
  }
} satisfies Actions;
