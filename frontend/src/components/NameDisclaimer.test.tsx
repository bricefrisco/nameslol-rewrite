import { expect, test } from "vitest";
import { render, screen } from "@testing-library/react";
import NameDisclaimer from "./NameDisclaimer.tsx";

test("renders information circle icon", () => {
  render(<NameDisclaimer />);
  expect(screen.getByTestId("information-circle-icon")).toBeInTheDocument();
});

test("renders a disclaimer", () => {
  render(<NameDisclaimer />);
  expect(
    screen.getByText(
      "The name could be invalid, blocked by Riot, or taken by a banned summoner.",
    ),
  ).toBeInTheDocument();
});

test("renders with correct classes", () => {
  render(<NameDisclaimer className="test12345" />);
  expect(screen.getByTestId("name-disclaimer")).toHaveClass("test12345");
});
