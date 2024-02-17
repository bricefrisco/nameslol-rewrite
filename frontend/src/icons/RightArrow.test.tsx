import { expect, test } from "vitest";
import { render, screen } from "@testing-library/react";
import RightArrow from "./RightArrow.tsx";

test("renders a svg", () => {
  render(<RightArrow />);
  expect(screen.getByTestId("right-arrow-icon")).toBeInTheDocument();
  expect(screen.getByTestId("right-arrow-icon").tagName).toBe("svg");
});

test("renders with specified class", () => {
  render(<RightArrow className="test-class" />);
  expect(screen.getByTestId("right-arrow-icon")).toHaveClass("test-class");
});
