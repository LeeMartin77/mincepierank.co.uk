import type { PageServerLoadEvent } from './$types';
import {
  addAverageScore,
  getAllMakerPies,
  getAllPieRankingSummaries,
  getLatestRanking,
  getMincePieMakersForYear
} from '$lib/storage';
import { mapPiesAndRankings } from '$components/utilities/mapPiesAndRankings';
import { getConfig } from '$lib/storage/config';

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export const load = async (_event: PageServerLoadEvent) => {
  const config = await getConfig();
  const activeYear = parseInt(config.activeYear);

  const [dataRes, piesRes, rankingSummariesRes, latestRankingRes] = await Promise.all([
    getMincePieMakersForYear(activeYear),
    getAllMakerPies(activeYear),
    getAllPieRankingSummaries(activeYear),
    getLatestRanking(activeYear)
  ]);

  const [data, pies, rankingSummaries, latestRanking] = [
    dataRes.unwrapOr([]),
    piesRes.unwrapOr([]),
    rankingSummariesRes.unwrapOr([]),
    latestRankingRes.unwrapOr(undefined)
  ];

  const { mappedRankings, mappedPies, rankingOrder } = mapPiesAndRankings(pies, rankingSummaries);
  const topPieId = rankingOrder.shift();
  const latestPie = latestRanking
    ? mappedPies[latestRanking.makerid + '-' + latestRanking.pieid]
    : null;
  return {
    makers: data,
    latestPie,
    latestRanking: latestRanking
      ? {
          ...addAverageScore(latestRanking),
          userid: null,
          notes: null
        }
      : null,
    topPie: topPieId ? mappedPies[topPieId] : null,
    topPieRanking: topPieId ? mappedRankings[topPieId] : null
  };
};
