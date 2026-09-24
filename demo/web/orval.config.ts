import { defineConfig } from "orval";

export default defineConfig({
  ixdtfDemo: {
    input: {
      target: "../api/openapi.yaml",
    },
    output: {
      target: "./src/generated/api/endpoints.ts",
      schemas: "./src/generated/api/model",
      client: "react-query",
      httpClient: "fetch",
      mode: "split",
      override: {
        operations: {
          // parse / roundtrip はリアルタイム解析 (F-2-1, F-3-1) のため
          // mutation でなく入力値をキーにした query として生成する
          parseIxdtf: {
            query: {
              useQuery: true,
            },
          },
          roundtripIxdtf: {
            query: {
              useQuery: true,
            },
          },
        },
      },
    },
  },
});
