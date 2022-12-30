import { render, screen } from "@testing-library/react";
//import userEvent from "@testing-library/user-event";
import "@testing-library/jest-dom";
import { SeasonBanner } from "./seasonBanner";

describe("PieSummaryLink", () => {
  test("Renders At All", async () => {
    render(<SeasonBanner />);
  });

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

  test.each(badDates)("Out of Season :: %s :: Info Message", (date) => {
    render(<SeasonBanner nowFn={() => date.getTime()} />);

    expect(screen.getByTestId("out-of-season-banner")).toBeInTheDocument();
  });
  test.each(midSeason)("In Season :: %s :: No Message", (date) => {
    render(<SeasonBanner nowFn={() => date.getTime()} />);

    expect(screen.queryByTestId("out-of-season-banner")).toBeNull();
    expect(screen.queryByTestId("near-season-end-banner")).toBeNull();
  });
  test("Close to end of season :: Warning message", () => {
    const seasonEndDate = new Date(2022, 11, 12);
    const nowDate = new Date(2022, 11, 9);
    render(
      <SeasonBanner
        nowFn={() => seasonEndDate.getTime()}
        nextEndDateFn={() => nowDate.getTime()}
      />
    );

    expect(screen.getByTestId("near-season-end-banner")).toBeInTheDocument();
  });
});
