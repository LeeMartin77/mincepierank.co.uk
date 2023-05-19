import { error } from '@sveltejs/kit';
import type { PageServerLoadEvent } from './$types';
import { addAverageScore, getAllMakerPies, getAllPieRankingSummaries, getLatestRanking, getMincePieMakers } from '$lib/storage';
import { mapPiesAndRankings } from '$components/mapPiesAndRankings';

// TODO: Make this truly data driven
const years = ['2022'];

export const load = (async ({ params }: PageServerLoadEvent) => {
    const { year } = params;
    if (!years.includes(year)) {
      throw error(404, 'Not Found')
    }
    const data = (await getMincePieMakers()).unwrapOr([]);
    const pies = (await getAllMakerPies(parseInt(year))).unwrapOr([]);
    const rankingSummaries = (await getAllPieRankingSummaries()).unwrapOr([]);
    const latestRanking = (await getLatestRanking()).unwrapOr(undefined);
    const { mappedRankings, mappedPies, rankingOrder } = mapPiesAndRankings(
      pies,
      rankingSummaries
    );
    const topPieId = rankingOrder.shift();
    const latestPie = latestRanking
      ? mappedPies[latestRanking.makerid + "-" + latestRanking.pieid]
      : null;
    return {
        makers: data,
        latestPie,
        latestRanking: latestRanking
          ? {
              ...addAverageScore(latestRanking),
              userid: null,
              notes: null,
            }
          : null,
        topPie: topPieId ? mappedPies[topPieId] : null,
        topPieRanking: topPieId ? mappedRankings[topPieId] : null,
    };
    
});