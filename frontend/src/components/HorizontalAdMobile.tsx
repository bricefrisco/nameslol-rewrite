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

const HorizontalAdMobile = ({ id, className }: Props) => {
  const [ad, setAd] = useState<any>();
  const location = useLocation();

  useEffect(() => {
    const createAd = async () => {
      const ad = await Promise.resolve(
        window.nitroAds.createAd(id, {
          demo: process.env.NEXT_PUBLIC_ENVIRONMENT !== "production",
          format: "display",
          sizes: [[320, 50]],
          mediaQuery: "(max-width: 777px)",
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
    ad?.onNavigate();
  }, [location]);

  return (
    <div
      id={id}
      className={`np-horizontal-mobile ${className ?? ""}`}
      data-testid="horizontal-ad-mobile"
    />
  );
};

export default HorizontalAdMobile;
