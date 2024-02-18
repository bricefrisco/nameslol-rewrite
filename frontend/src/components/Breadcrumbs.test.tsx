import { expect, test } from "vitest";
import { render, screen } from "@testing-library/react";
import Breadcrumbs from "./Breadcrumbs.tsx";
import { BrowserRouter } from "react-router-dom";

test("renders a nav element", () => {
  render(<Breadcrumbs breadcrumbs={[]} />);
  expect(screen.getByTestId("breadcrumbs")).toBeInTheDocument();
  expect(screen.getByTestId("breadcrumbs").tagName).toBe("NAV");
});

test("renders an ordered list", () => {
  render(<Breadcrumbs breadcrumbs={[]} />);
  expect(screen.getByTestId("breadcrumbs-ol")).toBeInTheDocument();
  expect(screen.getByTestId("breadcrumbs-ol").tagName).toBe("OL");
});

test("renders a link for each breadcrumb", () => {
  const breadcrumbs = [
    { name: "home", href: "/" },
    { name: "summoners", href: "/summoners" },
  ];
  render(
    <BrowserRouter>
      <Breadcrumbs breadcrumbs={breadcrumbs} />
    </BrowserRouter>,
  );
  expect(screen.getByTestId("home-link")).toBeInTheDocument();
  expect(screen.getByTestId("summoners-link")).toBeInTheDocument();
});

test("renders home icon if name='home'", () => {
  const breadcrumbs = [{ name: "home", href: "/" }];
  render(
    <BrowserRouter>
      <Breadcrumbs breadcrumbs={breadcrumbs} />
    </BrowserRouter>,
  );
  expect(screen.getByTestId("home-icon")).toBeInTheDocument();
});

test("renders name of breadcrumb in link", () => {
  const breadcrumbs = [{ name: "home", href: "/" }];
  render(
    <BrowserRouter>
      <Breadcrumbs breadcrumbs={breadcrumbs} />
    </BrowserRouter>,
  );
  expect(screen.getByText("home")).toBeInTheDocument();
});

test("sets href attribute of link to href prop", () => {
  const breadcrumbs = [{ name: "home", href: "/home-test" }];
  render(
    <BrowserRouter>
      <Breadcrumbs breadcrumbs={breadcrumbs} />
    </BrowserRouter>,
  );
  expect(screen.getByTestId("home-link").getAttribute("href")).toBe(
    "/home-test",
  );
});

test("renders right arrow if not last breadcrumb", () => {
  const breadcrumbs = [
    { name: "home", href: "/" },
    { name: "summoners", href: "/summoners" },
  ];
  render(
    <BrowserRouter>
      <Breadcrumbs breadcrumbs={breadcrumbs} />
    </BrowserRouter>,
  );
  expect(screen.getByTestId("right-arrow-icon")).toBeInTheDocument();
});

test("does not render right arrow if last breadcrumb", () => {
  const breadcrumbs = [{ name: "home", href: "/" }];
  render(
    <BrowserRouter>
      <Breadcrumbs breadcrumbs={breadcrumbs} />
    </BrowserRouter>,
  );
  expect(screen.queryByTestId("right-arrow-icon")).not.toBeInTheDocument();
});
