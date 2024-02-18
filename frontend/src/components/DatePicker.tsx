import { useState } from "react";
import Datepicker from "tailwind-datepicker-react";

interface Props {
  label: string;
  className?: string;
  onChange: (date: Date) => void;
  date: Date;
}

const DatePicker = ({ label, className, onChange, date }: Props) => {
  const [show, setShow] = useState(false);

  return (
    <div className={className} data-testid="datepicker-div">
      <label className="mb-2 block text-sm font-medium text-gray-300">
        {label}
      </label>

      <Datepicker
        key={date.toString()}
        options={{
          title: "Available on",
          autoHide: true,
          todayBtn: true,
          clearBtn: true,
          // theme: {
          //   background: "bg-gray-800",
          // },
          defaultDate: date,
        }}
        show={show}
        setShow={() => {
          setShow(!show);
        }}
        onChange={onChange}
      />
    </div>
  );
};

export default DatePicker;
