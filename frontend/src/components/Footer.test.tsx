import { expect, test } from "vitest";
import { render, screen } from "@testing-library/react";
import { BrowserRouter } from "react-router-dom";
import Footer from "./Footer.tsx";

test("renders a footer element", () => {
  render(
    <BrowserRouter>
      <Footer />
    </BrowserRouter>,
  );

  expect(screen.getByTestId("footer")).toBeInTheDocument();
  expect(screen.getByTestId("footer").tagName).toBe("FOOTER");
});

test("renders the current year", () => {
  render(
    <BrowserRouter>
      <Footer />
    </BrowserRouter>,
  );

  const year = new Date().getFullYear();
  const yearRegExp = new RegExp("Copyright.*" + year.toString());
  expect(screen.getByText(yearRegExp)).toBeInTheDocument();
});

test("links to the privacy policy", () => {
  render(
    <BrowserRouter>
      <Footer />
    </BrowserRouter>,
  );

  const privacyPolicy = screen.getByTestId("privacy-policy");
  expect(privacyPolicy).toBeInTheDocument();
  expect(privacyPolicy).toHaveAttribute("href", "/privacy-policy");
});

test("links to the faq", () => {
  render(
    <BrowserRouter>
      <Footer />
    </BrowserRouter>,
  );

  const faq = screen.getByTestId("faq");
  expect(faq).toBeInTheDocument();
  expect(faq).toHaveAttribute("href", "/faq");
});

test("links to the repo", () => {
  render(
    <BrowserRouter>
      <Footer />
    </BrowserRouter>,
  );

  const link = screen.getByTestId("github-link");
  expect(link).toBeInTheDocument();
  expect(link).toHaveAttribute(
    "href",
    "https://github.com/bricefrisco/NamesLoL",
  );
});

test("has ncmp", () => {
  render(
    <BrowserRouter>
      <Footer />
    </BrowserRouter>,
  );

  expect(screen.getByTestId("ncmp")).toBeInTheDocument();
  expect(screen.getByTestId("ncmp")).toHaveAttribute("id", "ncmp-consent-link");
});

test("has ccpa", () => {
  render(
    <BrowserRouter>
      <Footer />
    </BrowserRouter>,
  );

  expect(screen.getByTestId("ccpa")).toBeInTheDocument();
  expect(screen.getByTestId("ccpa")).toHaveAttribute("data-ccpa-link", "1");
});
