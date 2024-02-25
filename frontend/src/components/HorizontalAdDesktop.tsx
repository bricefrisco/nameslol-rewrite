/* eslint-disable
    @typescript-eslint/no-explicit-any,
    @typescript-eslint/no-unsafe-assignment,
    @typescript-eslint/no-unsafe-member-access,
    @typescript-eslint/no-unsafe-call
*/
import { useEffect, useState } from "react";
import { useLocation } from "react-router-dom";

interface Props {
  id: string;
  className?: string;
}

declare global {
  interface Window {
    nitroAds: {
      createAd: (id: string, options: any) => Promise<any>;
    };
  }
}

const HorizontalAdDesktop = ({ id, className }: Props) => {
  const [ad, setAd] = useState<any>();
  const location = useLocation();

  useEffect(() => {
    const createAd = async () => {
      const ad = await Promise.resolve(
        window.nitroAds.createAd(id, {
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
        }),
      );

      setAd(ad);
    };

    void createAd();
  }, [id]);

  useEffect(() => {
    if (!ad) return;
    ad.onNavigate();
  }, [location]);

  return (
    <div
      id={id}
      className={`np-horizontal-desktop ${className ?? ""}`}
      data-testid="horizontal-ad-desktop"
    />
  );
};

export default HorizontalAdDesktop;
