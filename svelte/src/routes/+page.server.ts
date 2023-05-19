import type { PageServerLoadEvent } from './$types';
import { addAverageScore, getAllMakerPies, getAllPieRankingSummaries, getLatestRanking, getMincePieMakers } from '$lib/storage';
import { mapPiesAndRankings } from '../components/mapPiesAndRankings';

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export const load = (async (_event: PageServerLoadEvent) => {
    const data = (await getMincePieMakers()).unwrapOr([]);
    const pies = (await getAllMakerPies(2022)).unwrapOr([]);
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