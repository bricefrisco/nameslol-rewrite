import Home from "../icons/Home";
import RightArrow from "../icons/RightArrow.tsx";
import { Link } from "react-router-dom";

interface BreadcrumbsProps {
  breadcrumbs: BreadcrumbProps[];
}

const Breadcrumbs = ({ breadcrumbs }: BreadcrumbsProps) => {
  return (
    <nav className="flex" data-testid="breadcrumbs">
      <ol
        className="inline-flex items-center space-x-1 text-gray-400 md:space-x-3"
        data-testid="breadcrumbs-ol"
      >
        {breadcrumbs.map((breadcrumb, index) => {
          return (
            <Breadcrumb
              key={index}
              name={breadcrumb.name}
              href={breadcrumb.href}
              last={index === breadcrumbs.length - 1}
            />
          );
        })}
      </ol>
    </nav>
  );
};

interface BreadcrumbProps {
  name: string;
  href: string;
  last?: boolean;
}

const Breadcrumb = ({ name, href, last }: BreadcrumbProps) => {
  return (
    <>
      <Link to={href} data-testid={`${name}-link`}>
        <li className="flex cursor-pointer items-center text-sm font-medium capitalize hover:text-white">
          {name === "home" && <Home className="mr-2 h-4 w-4" />}
          {name}
        </li>
      </Link>
      {!last && <RightArrow className="h-6 w-6" />}
    </>
  );
};

export default Breadcrumbs;
