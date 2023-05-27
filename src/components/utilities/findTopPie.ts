import type { MakerPie, PieRankingSummary } from '$lib/storage';

export function findTopPie(
  pies: MakerPie[],
  rankings: PieRankingSummary[]
): [MakerPie, PieRankingSummary] | undefined {
  if (!pies || pies.length === 0 || !rankings || rankings.length === 0) {
    return undefined;
  }
  const topRanking = rankings
    .map((r) => [r.filling + r.looks + r.pastry + r.topping + r.value, r.makerid, r.pieid])
    .sort(([val], [valb]) => (val as number) - (valb as number))
    .pop();
  if (!topRanking) {
    return undefined;
  }
  const [, makerid, pieid] = topRanking;
  const pie = pies.find((p) => p.makerid === makerid && p.id === pieid);
  const ranking = rankings.find((r) => r.makerid === makerid && r.pieid === pieid);
  if (!pie || !ranking) {
    return undefined;
  }
  return [pie, ranking];
}
