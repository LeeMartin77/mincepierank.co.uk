// TODO: Break this up to be every year when we're supporting
// multiple years
export const SEASON_START_DATE = new Date(2022, 9, 1, 0, 0, 1).getTime();
export const SEASON_END_DATE = new Date(2022, 11, 31, 23, 59, 59).getTime();

export function getNextSeasonEndDate() {
  return SEASON_END_DATE;
}

export function timeInSeason(timems: number) {
  return timems < SEASON_END_DATE && timems > SEASON_START_DATE;
}
