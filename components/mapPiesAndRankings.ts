import { MakerPie } from "../system/storage";
import { PieListRanking } from "./pieList/pieList";

export function mapPiesAndRankings(
    pies: MakerPie[],
    rankings: PieListRanking[]): {
      mappedRankings: {
        [key: string]: PieListRanking | undefined;
      },
      mappedPies: { [key: string]: MakerPie },
      rankingOrder: string[],
      unrankedPies: Set<string>
    } {
  
      const mappedRankings: {
        [key: string]: PieListRanking | undefined;
      } = {};
      const mappedPies: { [key: string]: MakerPie } = {};
      const rankingOrder: string[] = [];
    
      const unrankedPies: Set<string> = new Set();
      pies.forEach((pie) => {
        const uniqId = pie.makerid + "-" + pie.id;
        mappedPies[uniqId] = pie;
        unrankedPies.add(uniqId);
      });
      rankings
        .sort((a, b) => b.average - a.average)
        .forEach((ranking) => {
          const uniqId = ranking.makerid + "-" + ranking.pieid;
          rankingOrder.push(uniqId);
          unrankedPies.delete(uniqId);
          mappedRankings[uniqId] = ranking;
        });
        return { mappedRankings, mappedPies, rankingOrder, unrankedPies }
    }
  