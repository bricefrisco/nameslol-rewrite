import { Link } from "react-router-dom";

const Footer = () => {
  const year = new Date().getFullYear();

  return (
    <footer
      className="container mx-auto block max-w-screen-xl place-content-between items-start border-t border-gray-700 px-3.5 py-6 text-sm md:py-8"
      data-testid="footer"
    >
      <div className="place-content-between lg:flex">
        <div>
          <h2 className="text-lg font-semibold">NamesLoL</h2>
          <p className="mt-1 mb-2 text-sm text-gray-400">
            Copyright &copy; {year}.{" "}
            <Link
              className="text-blue-400 hover:underline"
              to="/privacy-policy"
              data-testid="privacy-policy"
            >
              Privacy Policy
            </Link>
          </p>
          <p className="mt-2 max-w-lg text-sm leading-5 text-gray-400">
            NamesLoL isn&apos;t endorsed by Riot Games and doesn&apos;t reflect
            the views or opinions of Riot Games or anyone officially involved in
            producing or managing Riot Games properties. Riot Games, and all
            associated properties are trademarks or registered trademarks of
            Riot Games, Inc.
          </p>

          {/* ?usp_debug=1 */}
          <div id="ncmp-consent-link" className="mt-2" data-testid="ncmp"></div>

          {/* ?gdpr_debug=1 */}
          <span
            data-ccpa-link="1"
            className="ccpa mt-2"
            data-testid="ccpa"
          ></span>
        </div>
        <div>
          <p className="max-w-lg text-sm leading-5 text-gray-400 lg:max-w-xs lg:text-right">
            Open source. Like this project?{" "}
            <a
              className="text-blue-400 hover:underline"
              rel="noopener noreferrer"
              target="_blank"
              href="https://github.com/bricefrisco/NamesLoL"
              data-testid="github-link"
            >
              View the code
            </a>{" "}
            and give it a star! To report a bug or request a feature, please{" "}
            <a
              className="text-blue-400 hover:underline"
              rel="noopener noreferrer"
              target="_blank"
              href="https://github.com/bricefrisco/NamesLoL/issues"
            >
              create an issue
            </a>
            .
          </p>
          <p className="mt-2 max-w-lg text-sm leading-5 text-gray-400 lg:max-w-xs lg:text-right">
            Visit our{" "}
            <Link
              to="/faq"
              className="text-blue-400 hover:underline"
              data-testid="faq"
            >
              FAQ
            </Link>{" "}
            to learn how NamesLoL works and discover tips on securing your next
            name.
          </p>
        </div>
      </div>
    </footer>
  );
};

export default Footer;
