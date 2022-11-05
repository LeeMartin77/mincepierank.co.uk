import { render, screen } from "@testing-library/react";
//import userEvent from "@testing-library/user-event";
import "@testing-library/jest-dom";
import { PieSummaryLink } from "./pieSummaryLink";
import { MakerPie } from "../system/storage";

describe("PieSummaryLink", () => {
  test("loads and displays pie name", async () => {
    // ARRANGE
    const pie: MakerPie = {
      makerid: "some-maker-id",
      id: "some-pie-id",
      displayname: "Some Pie Test",
      fresh: false,
      labels: [],
    };

    render(<PieSummaryLink pie={pie} />);

    // ACT
    //await userEvent.click(screen.getByText("Dummy Component"));
    //await screen.findByRole('heading')

    // ASSERT
    expect(screen.getByText(pie.displayname)).toHaveTextContent(
      pie.displayname
    );
    //expect(screen.getByRole('heading')).toHaveTextContent('hello there')
    //expect(screen.getByRole('button')).toBeDisabled()
  });

  test("should have link(s) to pie", () => {
    // ARRANGE
    const pie: MakerPie = {
      makerid: "some-maker-id",
      id: "some-pie-id",
      displayname: "Some Pie Test",
      fresh: false,
      labels: [],
    };

    // ACT
    render(<PieSummaryLink pie={pie} />);

    // ASSERT
    const links = screen.getAllByRole("link");
    expect(links.length).toBeGreaterThan(0);
    links.forEach((link) =>
      // TODO: Look into this localhost thing
      expect(link).toHaveProperty(
        "href",
        `http://localhost/brands/${pie.makerid}/${pie.id}`
      )
    );
  });
});
