import { expect, test } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";
import Header from "./Header.tsx";
import { BrowserRouter } from "react-router-dom";

const renderWithProps = ({
  navItems,
}: {
  navItems: { href: string; label: string }[];
}) => {
  return render(
    <BrowserRouter>
      <Header navItems={navItems} />
    </BrowserRouter>,
  );
};

test("renders a header element", () => {
  renderWithProps({ navItems: [] });
  expect(screen.getByTestId("header")).toBeInTheDocument();
  expect(screen.getByTestId("header").tagName).toBe("HEADER");
});

test("renders a logo", () => {
  renderWithProps({ navItems: [] });
  expect(screen.getByTestId("logo")).toBeInTheDocument();
  expect(screen.getByTestId("logo").tagName).toBe("IMG");
});

test("contains link to home", () => {
  renderWithProps({ navItems: [] });
  const links = screen.getAllByRole("link");
  expect(links).toHaveLength(1);
  expect(links[0]).toHaveAttribute("href", "/");
});

test("renders navigation items", () => {
  renderWithProps({
    navItems: [
      { href: "/test1", label: "Test 1" },
      { href: "/test2", label: "Test 2" },
    ],
  });

  expect(screen.getByText("Test 1")).toBeInTheDocument();
  expect(screen.getByText("Test 2")).toBeInTheDocument();

  const links = screen.getAllByRole("link");
  expect(links).toHaveLength(3);
  expect(links[1]).toHaveAttribute("href", "/test1");
  expect(links[2]).toHaveAttribute("href", "/test2");
});

test("renders hamburger button", () => {
  renderWithProps({ navItems: [] });
  expect(screen.getByTestId("navbar-hamburger")).toBeInTheDocument();
});

test("hamburger navigation disabled by default", () => {
  renderWithProps({ navItems: [] });
  expect(screen.getByTestId("navbar-hamburger-items")).toHaveClass("hidden");
});

test("hamburger navigation toggles", () => {
  renderWithProps({ navItems: [] });
  const button = screen.getByTestId("navbar-hamburger");
  fireEvent.click(button);
  expect(screen.getByTestId("navbar-hamburger-items")).not.toHaveClass(
    "hidden",
  );
});
