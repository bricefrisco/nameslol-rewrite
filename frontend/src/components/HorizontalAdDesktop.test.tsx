import { expect, test, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import HorizontalAdDesktop from "./HorizontalAdDesktop.tsx";
import { act } from "react-dom/test-utils";
import { BrowserRouter, useLocation } from "react-router-dom";

vi.mock("react-router-dom", async (importOriginal) => {
  const mod = await importOriginal<typeof import("react-router-dom")>();
  return {
    ...mod,
    useLocation: vi.fn(),
  };
});

const renderWithProps = async ({ id }: { id: string }) => {
  return act(() => {
    return render(
      <BrowserRouter>
        <HorizontalAdDesktop id={id} />,
      </BrowserRouter>,
    );
  });
};

const onNavigate = vi.fn();
const createAd = vi.fn().mockResolvedValue({ onNavigate });

beforeEach(() => {
  vi.mocked(useLocation).mockClear();
  onNavigate.mockClear();
  vi.stubGlobal("nitroAds", {
    createAd,
  });
});

test("renders a div element", async () => {
  await renderWithProps({ id: "test" });
  expect(screen.getByTestId("horizontal-ad-desktop")).toBeInTheDocument();
  expect(screen.getByTestId("horizontal-ad-desktop").tagName).toBe("DIV");
});

test("calls nitroAds.createAd with correct parameters", async () => {
  await renderWithProps({ id: "test" });
  expect(createAd).toHaveBeenCalledWith("test", {
    demo: process.env.NODE_ENV !== "production",
    format: "display",
    sizes: [[728, 90]],
    mediaQuery: "(min-width: 778px)",
    refreshVisibleOnly: true,
    renderVisibleOnly: true,
    refreshLimit: 10,
    refreshTime: 60,
    report: {
      enabled: true,
    },
  });
});

test("when query parameter changes, calls ad.onNavigate", async () => {
  const { rerender } = await renderWithProps({ id: "test" });
  expect(onNavigate).not.toHaveBeenCalled();

  vi.mocked(useLocation).mockReturnValue({
    hash: "",
    key: "",
    pathname: "",
    state: undefined,
    search: "?test=1",
  });

  act(() => {
    rerender(
      <BrowserRouter>
        <HorizontalAdDesktop id="test" />,
      </BrowserRouter>,
    );
  });

  expect(onNavigate).toHaveBeenCalled();
});
