import { getPiesWithCategory } from '$lib/storage';
import { error } from '@sveltejs/kit';
import type { PageServerLoadEvent } from './$types';

export const load = async ({ params }: PageServerLoadEvent) => {
	const { year, category } = params;
	const pies = await getPiesWithCategory(parseInt(year), category);
	if (pies.isErr()) {
		throw error(404, 'Not Found');
	}
	return {
		pies: pies.value
	};
};
