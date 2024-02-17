import { expect, test } from "vitest";
import { render, screen } from "@testing-library/react";
import CheckCircle from "./CheckCircle.tsx";

test("renders a svg", () => {
  render(<CheckCircle />);
  expect(screen.getByTestId("check-circle-icon")).toBeInTheDocument();
  expect(screen.getByTestId("check-circle-icon").tagName).toBe("svg");
});

test("renders with specified class", () => {
  render(<CheckCircle className="test-class" />);
  expect(screen.getByTestId("check-circle-icon")).toHaveClass("test-class");
});
