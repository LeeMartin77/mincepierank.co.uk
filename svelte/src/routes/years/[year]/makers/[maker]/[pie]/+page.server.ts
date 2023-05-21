import type { PageServerLoadEvent } from './$types';
import { getPieByMakerAndId } from '$lib/storage';
import { error } from '@sveltejs/kit';
export const load = async ({ params }: PageServerLoadEvent) => {
	const { year, maker, pie } = params;
	const res = await getPieByMakerAndId(parseInt(year), maker, pie);
	if (res.isErr()) {
		throw error(404, 'Not Found');
	}

	return {
		pie: res.value
	};
};
