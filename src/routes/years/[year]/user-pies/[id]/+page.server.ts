import type { PageServerLoadEvent } from './$types';
import { deleteUserPie, getUserPieByIdWithOwner, getUserPieRankingSummary } from '$lib/storage';
import { error, redirect, type Actions } from '@sveltejs/kit';
export const load = async ({ params, parent }: PageServerLoadEvent) => {
  const { year, id } = params;
  const res = await getUserPieByIdWithOwner(parseInt(year), id);
  if (res.isErr()) {
    throw error(404, 'Not Found');
  }
  const rankingRes = await getUserPieRankingSummary(parseInt(year), id);

  const { owner_userid, ...pie } = res.value;

  const { session } = await parent();
  const userid = session?.user?.email;
  return {
    isMine: !!owner_userid && owner_userid === userid,
    pie,
    ranking: rankingRes.isOk() ? rankingRes.value : undefined
  };
};

export const actions = {
  delete: async ({ request, locals }) => {
    const session = await locals.getSession();
    const userid = session?.user?.email;
    const data = await request.formData();
    const year = parseInt((data.get('year') as string) || '');
    const id = data.get('id') as string;
    const res = await getUserPieByIdWithOwner(year, id);
    if (res.isErr() || !res.value.owner_userid || res.value.owner_userid !== userid) {
      throw redirect(303, `/profile/pies/`);
    }

    await deleteUserPie(year, id);

    throw redirect(303, `/profile/pies/`);
  }
} satisfies Actions;
