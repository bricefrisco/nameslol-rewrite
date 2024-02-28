import { expect, test, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import { act } from "react-dom/test-utils";
import { BrowserRouter, useLocation } from "react-router-dom";
import SiderailAd from "./SiderailAd.tsx";

vi.mock("react-router-dom", async (importOriginal) => {
  const mod = await importOriginal<typeof import("react-router-dom")>();
  return {
    ...mod,
    useLocation: vi.fn(),
  };
});

const onNavigate = vi.fn();
const createAd = vi.fn().mockResolvedValue({ onNavigate });

vi.stubGlobal("nitroAds", {
  createAd,
});

const renderWithProps = ({
  id,
  mediaMinWidth,
  className,
}: {
  id: string;
  mediaMinWidth?: string;
  className?: string;
}) => {
  return act(() => {
    return render(
      <BrowserRouter>
        <SiderailAd
          id={id}
          mediaMinWidth={mediaMinWidth}
          className={className}
        />
        ,
      </BrowserRouter>,
    );
  });
};

beforeEach(() => {
  vi.mocked(useLocation).mockClear();
  onNavigate.mockClear();
});

test("renders a div element", async () => {
  await renderWithProps({ id: "test" });
  expect(screen.getByTestId("siderail-ad")).toBeInTheDocument();
  expect(screen.getByTestId("siderail-ad").tagName).toBe("DIV");
});

test("renders with classname", async () => {
  await renderWithProps({ id: "test", className: "test-class" });
  expect(screen.getByTestId("siderail-ad")).toHaveClass("test-class");
});

test("calls nitroAds.createAd with correct parameters", async () => {
  await renderWithProps({ id: "test" });
  expect(createAd).toHaveBeenCalledWith("test", {
    demo: process.env.NEXT_PUBLIC_ENVIRONMENT !== "production",
    format: "display",
    sizes: [[160, 600]],
    mediaQuery: "(min-width: 1600px) and (min-height: 700px)",
    refreshVisibleOnly: true,
    renderVisibleOnly: true,
    refreshLimit: 10,
    refreshTime: 60,
    report: {
      enabled: true,
    },
  });
});

test("calls nitroAds.createAd with correct parameters when mediaMinWidth is provided", async () => {
  await renderWithProps({ id: "test", mediaMinWidth: "800px" });
  expect(createAd).toHaveBeenCalledWith("test", {
    demo: process.env.NEXT_PUBLIC_ENVIRONMENT !== "production",
    format: "display",
    sizes: [[160, 600]],
    mediaQuery: "(min-width: 800px) and (min-height: 700px)",
    refreshVisibleOnly: true,
    renderVisibleOnly: true,
    refreshLimit: 10,
    refreshTime: 60,
    report: {
      enabled: true,
    },
  });
});

test("calls onNavigate when location changes", async () => {
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
        <SiderailAd id="test" />,
      </BrowserRouter>,
    );
  });

  expect(onNavigate).toHaveBeenCalledTimes(1);
});
