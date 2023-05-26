import {
	getPieRankingSummariesByIds,
	getPiesWithCategory,
	type PieRankingSummary
} from '$lib/storage';
import { error } from '@sveltejs/kit';
import type { PageServerLoad, PageServerLoadEvent } from './$types';
import { mapPiesAndRankings } from '$components/utilities/mapPiesAndRankings';

export const load = (async ({ params }: PageServerLoadEvent) => {
	const { year, category } = params;
	const pies = await getPiesWithCategory(parseInt(year), category);
	if (pies.isErr()) {
		throw error(404, 'Not Found');
	}
	const makerids = pies.value.reduce<{ [key: string]: string[] }>((acc, pie) => {
		if (acc[pie.makerid]) {
			acc[pie.makerid].push(pie.id);
		} else {
			acc[pie.makerid] = [pie.id];
		}
		return acc;
	}, {});

	const pieRankingResults = await Promise.all(
		Object.entries(makerids).map(([maker, pieids]) =>
			getPieRankingSummariesByIds(parseInt(year), maker, pieids)
		)
	);
	const pieRankings = pieRankingResults.reduce<PieRankingSummary[]>((acc, res) => {
		if (res.isOk()) {
			return [...acc, ...res.value];
		}
		return acc;
	}, []);
	return {
		pies: pies.value,
		pieRankings,
		category
	};
}) satisfies PageServerLoad;
