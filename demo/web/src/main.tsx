import { RouterProvider } from "@tanstack/react-router";
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import "./index.css";
import { Providers } from "@/app/providers";
import { router } from "@/app/router";
import { UnsupportedBrowser } from "@/components/UnsupportedBrowser";
import { isTemporalAvailable } from "@/lib/temporal/detect";

const rootElement = document.getElementById("root");
if (!rootElement) {
  throw new Error("#root element not found");
}

// F-0-4: ネイティブ Temporal の feature detection。ポリフィルは使わない (N-1)。
createRoot(rootElement).render(
  <StrictMode>
    <Providers>
      {isTemporalAvailable() ? <RouterProvider router={router} /> : <UnsupportedBrowser />}
    </Providers>
  </StrictMode>,
);
