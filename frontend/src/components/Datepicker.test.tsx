import { expect, test, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import { BrowserRouter } from "react-router-dom";
import DatePicker from "./DatePicker.tsx";

const onChange = vi.fn();

interface renderWithPropsProps {
  label?: string;
  className?: string;
  onChangeProp?: (date: Date) => void;
  date?: Date;
}

const renderWithProps = ({
  label,
  className,
  onChangeProp,
  date,
}: renderWithPropsProps) => {
  const props = {
    label: label ?? "Available on",
    className: className ?? "test",
    onChange: onChangeProp ?? onChange,
    date: date ?? new Date(),
  };

  render(
    <BrowserRouter>
      <DatePicker {...props} />
    </BrowserRouter>,
  );
};

afterEach(() => {
  vi.resetAllMocks();
});

test("renders a datepicker", () => {
  renderWithProps({});
  expect(screen.getByTestId("datepicker-div")).toBeInTheDocument();
});

test("sets classnames properly", () => {
  renderWithProps({ className: "test-123" });
  expect(screen.getByTestId("datepicker-div")).toHaveClass("test-123");
});

test("displays label properly", () => {
  renderWithProps({ label: "Test label" });
  expect(screen.getByText("Test label")).toBeInTheDocument();
});
