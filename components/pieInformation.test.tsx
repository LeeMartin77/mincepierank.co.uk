import { render, screen } from "@testing-library/react";
//import userEvent from "@testing-library/user-event";
import "@testing-library/jest-dom";
import { PieSummaryInformation } from "./pieInformation";
import { MakerPie } from "../system/storage";

test("loads and displays pie name", async () => {
  // ARRANGE
  const pie: MakerPie = {
    makerid: "some-maker-id",
    id: "some-pie-id",
    displayname: "Some Pie Test",
    fresh: false,
    labels: [],
  };

  render(<PieSummaryInformation pie={pie} />);

  // ACT
  //await userEvent.click(screen.getByText("Dummy Component"));
  //await screen.findByRole('heading')

  // ASSERT
  expect(screen.getByText(pie.displayname)).toHaveTextContent(pie.displayname);
  //expect(screen.getByRole('heading')).toHaveTextContent('hello there')
  //expect(screen.getByRole('button')).toBeDisabled()
});
