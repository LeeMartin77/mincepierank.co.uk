import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import "@testing-library/jest-dom";
import { Dummy } from "./dummy";

test("loads and displays greeting", async () => {
  // ARRANGE
  render(<Dummy />);

  // ACT
  await userEvent.click(screen.getByText("Dummy Component"));
  //await screen.findByRole('heading')

  // ASSERT
  expect(screen.getByText("Dummy Component")).toHaveTextContent(
    "Dummy Component"
  );
  //expect(screen.getByRole('heading')).toHaveTextContent('hello there')
  //expect(screen.getByRole('button')).toBeDisabled()
});
