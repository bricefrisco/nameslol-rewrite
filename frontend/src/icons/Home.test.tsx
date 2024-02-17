import { expect, test } from "vitest";
import { render, screen } from "@testing-library/react";
import Home from "./Home.tsx";

test("renders a svg", () => {
  render(<Home />);
  expect(screen.getByTestId("home-icon")).toBeInTheDocument();
  expect(screen.getByTestId("home-icon").tagName).toBe("svg");
});

test("renders with specified class", () => {
  render(<Home className="test-class" />);
  expect(screen.getByTestId("home-icon")).toHaveClass("test-class");
});
