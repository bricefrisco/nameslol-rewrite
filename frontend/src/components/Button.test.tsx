import { expect, test } from "vitest";
import { render, screen } from "@testing-library/react";
import { BrowserRouter } from "react-router-dom";
import Button from "./Button.tsx";

test("renders a button when href is not provided", () => {
  render(
    <BrowserRouter>
      <Button>Click me</Button>
    </BrowserRouter>,
  );
  expect(screen.getByTestId("button-button")).toBeInTheDocument();
});

test("renders a link when href is provided", () => {
  render(
    <BrowserRouter>
      <Button href="/test">Click me</Button>
    </BrowserRouter>,
  );
  expect(screen.getByTestId("button-link")).toBeInTheDocument();
});

test("sets classnames on button", () => {
  render(
    <BrowserRouter>
      <Button className="test">Click me</Button>
    </BrowserRouter>,
  );
  expect(screen.getByTestId("button-div").className).toBe("test");
});

test("sets classnames on a link", () => {
  render(
    <BrowserRouter>
      <Button href="/test" className="test">
        Click me
      </Button>
    </BrowserRouter>,
  );
  expect(screen.getByTestId("button-div").className).toBe("test");
});

test("sets type attribute on button", () => {
  render(
    <BrowserRouter>
      <Button type="submit">Click me</Button>
    </BrowserRouter>,
  );
  expect(screen.getByTestId("button-button").getAttribute("type")).toBe(
    "submit",
  );
});

it("renders children (button)", () => {
  render(
    <BrowserRouter>
      <Button>Click me</Button>
    </BrowserRouter>,
  );
  expect(screen.getByText("Click me")).toBeInTheDocument();
});

it("renders children (link)", () => {
  render(
    <BrowserRouter>
      <Button href="/test">Click me</Button>
    </BrowserRouter>,
  );
  expect(screen.getByText("Click me")).toBeInTheDocument();
});

it("sets href attribute of link to href prop", () => {
  render(
    <BrowserRouter>
      <Button href="/test">Click me</Button>
    </BrowserRouter>,
  );
  expect(screen.getByTestId("button-link").getAttribute("href")).toBe("/test");
});
