import { createRootRoute, createRoute, createRouter, redirect } from "@tanstack/react-router";
import { Layout } from "@/components/Layout";
import { AboutPage } from "@/pages/about/AboutPage";
import { ConverterPage } from "@/pages/converter/ConverterPage";
import { GuidePage } from "@/pages/guide/GuidePage";
import { HomePage } from "@/pages/home/HomePage";
import { SupportPage } from "@/pages/support/SupportPage";
import { WorkbenchPage } from "@/pages/workbench/WorkbenchPage";

/** ワークベンチの表示モード (F-2 解析・検証 / F-3 往復・実装差比較 / F-6 Temporal ラボ)。 */
export type WorkbenchMode = "parse" | "roundtrip" | "lab";

/** Typed search params for the Workbench shareable URL (F-2-7). */
export interface PlaygroundSearch {
  input?: string;
  strict?: boolean;
  validateOnly?: boolean;
  mode?: WorkbenchMode;
}

const rootRoute = createRootRoute({ component: Layout });

const homeRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/",
  component: HomePage,
});

const aboutRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/about",
  component: AboutPage,
});

const playgroundRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/playground",
  component: WorkbenchPage,
  validateSearch: (search): PlaygroundSearch => ({
    input: typeof search.input === "string" ? search.input : undefined,
    strict: search.strict === true || search.strict === "true" ? true : undefined,
    validateOnly: search.validateOnly === true || search.validateOnly === "true" ? true : undefined,
    mode: search.mode === "roundtrip" || search.mode === "lab" ? search.mode : undefined,
  }),
});

// F-3 は F-2 ワークベンチの 1 モードへ格下げ。旧 /interop の deep link は
// /playground?mode=roundtrip へリダイレクトして保持する。
const interopRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/interop",
  beforeLoad: () => {
    throw redirect({ to: "/playground", search: { mode: "roundtrip" }, replace: true });
  },
});

const converterRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/converter",
  component: ConverterPage,
});

const guideRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/guide",
  component: GuidePage,
});

// F-6 は F-2 ワークベンチの Lab モードへ格下げ (D-12 反転)。旧 /temporal-lab の
// deep link は /playground?mode=lab へリダイレクトして保持する。
const temporalLabRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/temporal-lab",
  beforeLoad: () => {
    throw redirect({ to: "/playground", search: { mode: "lab" }, replace: true });
  },
});

const supportRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/support",
  component: SupportPage,
});

const routeTree = rootRoute.addChildren([
  homeRoute,
  aboutRoute,
  playgroundRoute,
  interopRoute,
  converterRoute,
  guideRoute,
  temporalLabRoute,
  supportRoute,
]);

export const router = createRouter({ routeTree, defaultPreload: "intent" });

declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router;
  }
}
