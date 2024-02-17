import { expect, test } from "vitest";
import { render, screen } from "@testing-library/react";
import BookOpen from "./BookOpen.tsx";

test("testing", () => {
  render(<BookOpen />);
  expect(screen.getByTestId("book-open-icon")).toBeInTheDocument();
});
