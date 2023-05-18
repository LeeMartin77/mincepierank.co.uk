import { timeInSeason } from "./seasonConfig";

describe("Season Config", () => {
  const badDates = [
    new Date(2023, 0, 1, 0, 1, 0),
    new Date(2022, 8, 30, 22, 58, 0),
    new Date(2022, 4, 12, 12, 30, 0),
  ];

  const midSeason = [
    new Date(2022, 11, 12, 12, 30, 0),
    new Date(2022, 10, 12, 12, 30, 0),
    new Date(2022, 9, 12, 12, 30, 0),
  ];

  test.each(badDates)("%s Out of season", (date) => {
    expect(timeInSeason(date.getTime())).toBeFalse();
  });
  test.each(midSeason)("%s In season", (date) => {
    expect(timeInSeason(date.getTime())).toBeTrue();
  });
});
