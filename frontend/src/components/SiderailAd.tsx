/* eslint-disable
    @typescript-eslint/no-explicit-any,
    @typescript-eslint/no-unsafe-assignment,
    @typescript-eslint/no-unsafe-member-access,
    @typescript-eslint/no-unsafe-call
*/
import { useState, useEffect } from "react";
import { useLocation } from "react-router-dom";

interface Props {
  id: string;
  mediaMinWidth?: string;
  className?: string;
}

const SideRailAd = ({ id, mediaMinWidth, className }: Props) => {
  const [ad, setAd] = useState<any>();
  const location = useLocation();

  useEffect(() => {
    const createAd = async () => {
      const ad = await Promise.resolve(
        window.nitroAds.createAd(id, {
          demo: process.env.NEXT_PUBLIC_ENVIRONMENT !== "production",
          format: "display",
          sizes: [[160, 600]],
          mediaQuery: `(min-width: ${
            mediaMinWidth ?? "1600px"
          }) and (min-height: 700px)`,
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
  }, [id, mediaMinWidth]);

  useEffect(() => {
    ad?.onNavigate();
  }, [location]);

  return (
    <div
      id={id}
      className={`np-siderail ${className ?? ""}`}
      data-testid="siderail-ad"
    />
  );
};

export default SideRailAd;
