import { error } from '@sveltejs/kit';
import type { PageServerLoadEvent } from './$types';
import { getMincePieMakersForYear } from '$lib/storage';

// TODO: Make this truly data driven
const years = ['2022'];

export const load = async ({ params }: PageServerLoadEvent) => {
	const { year } = params;
	if (!years.includes(year)) {
		throw error(404, 'Not Found');
	}
	const yearInt = parseInt(year);
	const data = await getMincePieMakersForYear(yearInt);
	return {
		makers: data.isOk() ? data.value : [],
		year: yearInt
	};
};
