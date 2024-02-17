import { expect, test } from "vitest";
import { render, screen } from "@testing-library/react";
import BookOpen from "./BookOpen.tsx";

test("renders an svg", () => {
  render(<BookOpen />);
  expect(screen.getByTestId("book-open-icon")).toBeInTheDocument();
  expect(screen.getByTestId("book-open-icon").tagName).toBe("svg");
});

test("renders with specified class", () => {
  render(<BookOpen className="test-class" />);
  expect(screen.getByTestId("book-open-icon")).toHaveClass("test-class");
});
