import React from "react";
import { Link } from "react-router-dom";

interface ButtonProps {
  href?: string;
  className?: string;
  type?: "submit" | "reset" | "button" | undefined;
  children: React.ReactNode;
}

const Button = ({ href, className, type, children }: ButtonProps) => {
  const classes =
    "mr-2 mb-2 rounded-lg bg-blue-600 px-5 py-2.5 text-sm font-medium text-white hover:bg-blue-700 focus:outline-none focus:ring-4 focus:ring-blue-800";

  if (href) {
    return (
      <div className={className} data-testid="button-div">
        <Link to={href} data-testid="button-link">
          <span className={classes}>{children}</span>
        </Link>
      </div>
    );
  }

  return (
    <div className={className} data-testid="button-div">
      <button type={type} className={classes} data-testid="button-button">
        {children}
      </button>
    </div>
  );
};

export default Button;
