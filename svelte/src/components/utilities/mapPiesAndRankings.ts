import type { MakerPie, MakerPieRanking } from '$lib/storage';

export type PieListRanking = Omit<MakerPieRanking, 'userid' | 'notes'> & {
	count: number | undefined;
	average: number;
};

export function mapPiesAndRankings(
	pies: MakerPie[],
	rankings: PieListRanking[],
	filteredCategoryIds?: Set<string>
): {
	mappedRankings: {
		[key: string]: PieListRanking | undefined;
	};
	mappedPies: { [key: string]: MakerPie };
	rankingOrder: string[];
	unrankedPies: Set<string>;
} {
	const mappedRankings: {
		[key: string]: PieListRanking | undefined;
	} = {};
	const mappedPies: { [key: string]: MakerPie } = {};
	const rankingOrder: string[] = [];

	const unrankedPies: Set<string> = new Set();
	pies.forEach((pie) => {
		const uniqId = pie.makerid + '-' + pie.id;
		mappedPies[uniqId] = pie;
		if (
			filteredCategoryIds === undefined ||
			filteredCategoryIds.size === 0 ||
			Array.from(filteredCategoryIds).every((lb) => pie.labels.includes(lb))
		) {
			unrankedPies.add(uniqId);
		}
	});
	rankings
		.sort((a, b) => b.average - a.average)
		.forEach((ranking) => {
			const uniqId = ranking.makerid + '-' + ranking.pieid;
			mappedRankings[uniqId] = ranking;
			if (unrankedPies.has(uniqId)) {
				rankingOrder.push(uniqId);
				unrankedPies.delete(uniqId);
			}
		});
	return { mappedRankings, mappedPies, rankingOrder, unrankedPies };
}
