import type { Decorator, Preview } from "@storybook/react-vite";
import "../src/index.css";

// DESIGN.md Storybook 規則: ライト/ダーク両テーマの toolbar 切替。
// アプリ本体と同じく <html> のクラスでトークンを切り替える (dark: バリアント不使用)。
const withTheme: Decorator = (Story, context) => {
  const theme = context.globals.theme === "dark" ? "dark" : "light";
  document.documentElement.classList.remove("light", "dark");
  document.documentElement.classList.add(theme);
  document.body.style.backgroundColor = "var(--background)";
  return Story();
};

const preview: Preview = {
  decorators: [withTheme],
  globalTypes: {
    theme: {
      description: "ライト/ダークのトークン切替 (DESIGN.md)",
      toolbar: {
        title: "Theme",
        icon: "circlehollow",
        items: ["light", "dark"],
        dynamicTitle: true,
      },
    },
  },
  initialGlobals: { theme: "light" },
};

export default preview;
