import { expect, test } from "vitest";
import { render, screen } from "@testing-library/react";
import InformationCircle from "./InformationCircle.tsx";

test("renders a svg", () => {
  render(<InformationCircle />);
  expect(screen.getByTestId("information-circle-icon")).toBeInTheDocument();
  expect(screen.getByTestId("information-circle-icon").tagName).toBe("svg");
});

test("renders with specified class", () => {
  render(<InformationCircle className="test-class" />);
  expect(screen.getByTestId("information-circle-icon")).toHaveClass(
    "test-class",
  );
});
