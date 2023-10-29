import type { MakerPie, MakerPieRanking, UserPie, UserPieRanking } from '$lib/storage';

export type PieListRanking =
  | (Omit<MakerPieRanking, 'userid' | 'notes'> & {
      count: number | undefined;
      average: number;
    })
  | (Omit<UserPieRanking, 'userid' | 'notes'> & {
      count: number | undefined;
      average: number;
    });

export type MappedPiesAndRankings = {
  mappedRankings: {
    [key: string]: PieListRanking | undefined;
  };
  mappedPies: {
    [key: string]: UserPie | MakerPie;
  };
  rankingOrder: string[];
  unrankedPies: Set<string>;
};

export function mapPiesAndRankings(
  pies: (MakerPie | UserPie)[],
  rankings: PieListRanking[],
  filteredCategoryIds?: Set<string>
): MappedPiesAndRankings {
  const mappedRankings: {
    [key: string]: PieListRanking | undefined;
  } = {};
  const mappedPies: { [key: string]: MakerPie | UserPie } = {};
  const rankingOrder: string[] = [];

  const unrankedPies: Set<string> = new Set();
  pies.forEach((pie) => {
    let uniqId = pie.id;
    if ('makerid' in pie) {
      uniqId = pie.makerid + '-' + pie.id;
    }
    mappedPies[uniqId] = pie;
    if (
      filteredCategoryIds === undefined ||
      filteredCategoryIds.size === 0 ||
      Array.from(filteredCategoryIds).every((lb) => pie.labels && pie.labels.includes(lb))
    ) {
      unrankedPies.add(uniqId);
    }
  });
  rankings
    .sort((a, b) => b.average - a.average)
    .forEach((ranking) => {
      let uniqId = ranking.pieid;
      if ('makerid' in ranking) {
        uniqId = ranking.makerid + '-' + ranking.pieid;
      }
      mappedRankings[uniqId] = ranking;
      if (unrankedPies.has(uniqId)) {
        rankingOrder.push(uniqId);
        unrankedPies.delete(uniqId);
      }
    });
  return { mappedRankings, mappedPies, rankingOrder, unrankedPies };
}
