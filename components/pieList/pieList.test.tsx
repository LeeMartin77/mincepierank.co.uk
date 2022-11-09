import { PieList } from "./pieList";
import { render, screen } from "@testing-library/react";
//import userEvent from "@testing-library/user-event";
import "@testing-library/jest-dom";
import { MakerPie } from "../../system/storage";

describe("PieSummaryLink", () => {
  test("Renders At All", async () => {
    render(<PieList pies={[]} rankings={[]} />);
  });

  test("Has no rankings :: Just displays all pies in a list in original order", async () => {
    // ARRANGE
    const pies: MakerPie[] = [
      {
        makerid: "some-maker-id-1",
        id: "some-pie-id-1",
        displayname: "Some Pie Test 1",
        fresh: false,
        labels: [],
        image_file: "/image/pie1.jpg",
        web_link: "",
        pack_count: 6,
      },
      {
        makerid: "some-maker-id-2",
        id: "some-pie-id-2",
        displayname: "Some Pie Test 2",
        fresh: false,
        labels: [],
        image_file: "/image/pie2.jpg",
        web_link: "",
        pack_count: 6,
      },
      {
        makerid: "some-maker-id-3",
        id: "some-pie-id-3",
        displayname: "Some Pie Test 3",
        fresh: false,
        labels: [],
        image_file: "/image/pie3.jpg",
        web_link: "",
        pack_count: 6,
      },
    ];

    render(<PieList pies={pies} rankings={[]} />);

    // ASSERT
    for (let pie of pies) {
      // expect(screen.getByText(pie.displayname)).toBeInTheDocument();
      const summary = screen.getByTestId(pie.id + "-link-card");
      expect(summary).not.toBeUndefined();
      expect(summary).toHaveTextContent(pie.displayname);
      const links = summary.getElementsByTagName("a");
      for (let link of links) {
        expect(link.href).toBe(
          `http://localhost/brands/${pie.makerid}/${pie.id}`
        );
      }
    }
  });
});
