import i18n from "i18next";
import { initReactI18next } from "react-i18next";
import en from "@/locales/en/translation.json";
import ja from "@/locales/ja/translation.json";

export type Language = "ja" | "en";

const STORAGE_KEY = "lang";

function initialLanguage(): Language {
  const stored = localStorage.getItem(STORAGE_KEY);
  if (stored === "ja" || stored === "en") {
    return stored;
  }
  return typeof navigator !== "undefined" && navigator.language.startsWith("ja") ? "ja" : "en";
}

void i18n.use(initReactI18next).init({
  resources: {
    ja: { translation: ja },
    en: { translation: en },
  },
  lng: initialLanguage(),
  fallbackLng: "en",
  interpolation: { escapeValue: false },
});

/** Change the UI language and persist the choice (F-0-2). */
export function setLanguage(language: Language): void {
  localStorage.setItem(STORAGE_KEY, language);
  void i18n.changeLanguage(language);
}

export default i18n;
