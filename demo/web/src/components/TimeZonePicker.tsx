import { useTranslation } from "react-i18next";
import { Combobox } from "@/components/ui/Combobox";
import { supportedTimeZones } from "@/lib/intl";

interface TimeZonePickerProps {
  value: string;
  onChange: (timeZone: string) => void;
  className?: string;
}

/** Searchable IANA time-zone picker (F-1-3, `Intl.supportedValuesOf`). */
export function TimeZonePicker({ value, onChange, className }: TimeZonePickerProps) {
  const { t } = useTranslation();
  return (
    <Combobox
      value={value}
      options={supportedTimeZones()}
      onChange={onChange}
      aria-label={t("home.timeZoneLabel")}
      placeholder="Asia/Tokyo"
      className={className}
    />
  );
}
