import React from "react";

interface InputProps extends React.InputHTMLAttributes<HTMLInputElement> {
  label: string;
  className: string;
  value: string;
  setValue: (value: string) => void;
}

const Input = ({ label, className, value, setValue, ...rest }: InputProps) => {
  return (
    <div className={className} data-testid="input-div">
      <label className="mb-2 block text-sm font-medium text-gray-300">
        {label}
      </label>
      <input
        {...rest}
        data-testid="input"
        type="text"
        className="focus:ring-primary-500 focus:order-primary-500 shadow-sm-light block w-full rounded-lg border border-gray-600 bg-gray-700 p-2.5 text-sm text-white placeholder-gray-400 shadow-sm"
        placeholder="Example"
        value={value || ""}
        onChange={(e) => {
          setValue(e.target.value);
        }}
        required
      />
    </div>
  );
};

export default Input;
