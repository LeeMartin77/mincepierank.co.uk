import { Maker, MakerPie } from "../system/storage";
import { PieListRanking } from "./pieList/pieList";

export function descriptionSummary(pie: MakerPie, pieListRanking: PieListRanking, maker?: Maker): string { 
    return `The current leader is ${pie.displayname} ${maker ? `from ${maker.name} `:``}with ${pieListRanking?.average.toFixed(1)} stars across ${pieListRanking?.count} votes.` 
}