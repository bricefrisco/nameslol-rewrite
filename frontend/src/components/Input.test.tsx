import { expect, test, vi } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";
import Input from "./Input.tsx";

interface LabelProps {
  label?: string;
  className?: string;
  value?: string;
  setValue?: (value: string) => void;
}

const renderWithProps = ({ label, className, value, setValue }: LabelProps) => {
  return render(
    <Input
      label={label ?? "test"}
      className={className ?? "test"}
      value={value ?? "test"}
      // eslint-disable-next-line @typescript-eslint/no-empty-function
      setValue={setValue ?? (() => {})}
    />,
  );
};

test("renders an input element", () => {
  renderWithProps({});
  expect(screen.getByTestId("input")).toBeInTheDocument();
  expect(screen.getByTestId("input").tagName).toBe("INPUT");
});

test("renders with correct classes", () => {
  renderWithProps({ className: "test12345" });
  expect(screen.getByTestId("input-div")).toHaveClass("test12345");
});

test("renders with a label", () => {
  renderWithProps({ label: "test label 123" });
  expect(screen.getByText("test label 123")).toBeInTheDocument();
});

test("renders with a value", () => {
  renderWithProps({ value: "test value 123" });
  expect(screen.getByDisplayValue("test value 123")).toBeInTheDocument();
});

test("calls setValue on change", () => {
  const setValue = vi.fn();
  renderWithProps({ setValue });
  const input = screen.getByTestId("input");
  fireEvent.change(input, { target: { value: "test value 123" } });
  input.dispatchEvent(new Event("input", { bubbles: true }));
  expect(setValue).toHaveBeenCalledWith("test value 123");
});

test("renders other props", () => {
  render(
    <Input
      label="test"
      className="test"
      value="test"
      // eslint-disable-next-line @typescript-eslint/no-empty-function
      setValue={() => {}}
      placeholder="test placeholder"
      data-test123="testing123"
    />,
  );

  const input = screen.getByTestId("input");
  expect(input).toHaveAttribute("data-test123", "testing123");
});

test("renders as required", () => {
  renderWithProps({});
  const input = screen.getByTestId("input");
  expect(input).toBeRequired();
});
