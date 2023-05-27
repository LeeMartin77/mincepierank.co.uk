import { getAllMakerPies, getPieRankingSummariesByIds, type PieRankingSummary } from '$lib/storage';
import { error } from '@sveltejs/kit';
import type { PageServerLoadEvent } from './$types';

export const load = async ({ params }: PageServerLoadEvent) => {
	const { year } = params;
	const pies = await getAllMakerPies(parseInt(year));
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
		pieRankings,
		categories,
		year: parseInt(year)
	};
};
