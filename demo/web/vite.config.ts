/// <reference types="vitest/config" />

import { fileURLToPath, URL } from "node:url";
import { storybookTest } from "@storybook/addon-vitest/vitest-plugin";
import tailwindcss from "@tailwindcss/vite";
import react from "@vitejs/plugin-react";
import { playwright } from "@vitest/browser-playwright";
import { defineConfig } from "vite";

// https://vite.dev/config/
export default defineConfig({
  plugins: [react(), tailwindcss()],
  resolve: {
    alias: {
      "@": fileURLToPath(new URL("./src", import.meta.url)),
    },
  },
  server: {
    // testdata/ixdtf_samples.json lives one level above web/ (shared fixture).
    fs: {
      allow: [".."],
    },
    proxy: {
      // 開発時は Go サーバー (:8080) へ /api をプロキシする (D-7)
      "/api": "http://localhost:8080",
    },
  },
  test: {
    projects: [
      {
        // 純ロジックの単体テスト (jsdom)。既存の *.test.ts(x) を実行する。
        extends: true,
        test: {
          name: "unit",
          environment: "jsdom",
          setupFiles: ["./vitest.setup.ts"],
          include: ["src/**/*.{test,spec}.{ts,tsx}"],
        },
      },
      {
        // Storybook の全 stories を実ブラウザ(chromium)で render/play テストとして実行する。
        extends: true,
        plugins: [storybookTest({ configDir: ".storybook" })],
        test: {
          name: "storybook",
          browser: {
            enabled: true,
            headless: true,
            provider: playwright({}),
            instances: [{ browser: "chromium" }],
          },
        },
      },
    ],
  },
});
