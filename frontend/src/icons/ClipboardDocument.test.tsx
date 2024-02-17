import { expect, test } from "vitest";
import { render, screen } from "@testing-library/react";
import ClipboardDocument from "./ClipboardDocument.tsx";

test("renders a svg", () => {
  render(<ClipboardDocument />);
  expect(screen.getByTestId("clipboard-document-icon")).toBeInTheDocument();
  expect(screen.getByTestId("clipboard-document-icon").tagName).toBe("svg");
});

test("renders with specified class", () => {
  render(<ClipboardDocument className="test-class" />);
  expect(screen.getByTestId("clipboard-document-icon")).toHaveClass(
    "test-class",
  );
});
