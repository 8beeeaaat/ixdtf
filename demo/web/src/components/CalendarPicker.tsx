import { useTranslation } from "react-i18next";
import { Combobox } from "@/components/ui/Combobox";
import { supportedCalendars } from "@/lib/intl";

interface CalendarPickerProps {
  value: string;
  onChange: (calendar: string) => void;
  /** FormField の label (htmlFor) と関連付けるための input id。 */
  id?: string;
  className?: string;
}

/** Calendar-system picker (F-1-4, `Intl.supportedValuesOf`). */
export function CalendarPicker({ value, onChange, id, className }: CalendarPickerProps) {
  const { t } = useTranslation();
  return (
    <Combobox
      value={value}
      options={supportedCalendars()}
      onChange={onChange}
      aria-label={t("home.calendarLabel")}
      placeholder="iso8601"
      id={id}
      className={className}
    />
  );
}
