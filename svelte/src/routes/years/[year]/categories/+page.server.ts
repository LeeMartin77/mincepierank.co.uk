import { getAllMakerPies } from '$lib/storage';
import { error } from '@sveltejs/kit';
import type { PageServerLoadEvent } from './$types';

export const load = async ({ params }: PageServerLoadEvent) => {
	const { year } = params;
	const pies = await getAllMakerPies(parseInt(year));
	if (pies.isErr()) {
		throw error(404, 'Not Found');
	}
	const categories = Array.from(
		pies.value.reduce<Set<string>>((acc, curr) => {
			curr.labels.forEach((label) => {
				acc.add(label);
			});
			return acc;
		}, new Set<string>())
	);
	return {
		pies: pies.value,
		categories,
		year: parseInt(year)
	};
};
