import { createRootRoute, createRoute, createRouter } from "@tanstack/react-router";
import { Layout } from "@/components/Layout";
import { AboutPage } from "@/pages/about/AboutPage";
import { ConverterPage } from "@/pages/converter/ConverterPage";
import { GuidePage } from "@/pages/guide/GuidePage";
import { HomePage } from "@/pages/home/HomePage";
import { InteropPage } from "@/pages/interop/InteropPage";
import { PlaygroundPage } from "@/pages/playground/PlaygroundPage";
import { SupportPage } from "@/pages/support/SupportPage";
import { TemporalLabPage } from "@/pages/temporal-lab/TemporalLabPage";

/** Typed search params for the Playground shareable URL (F-2-7). */
export interface PlaygroundSearch {
  input?: string;
  strict?: boolean;
  validateOnly?: boolean;
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
  component: PlaygroundPage,
  validateSearch: (search): PlaygroundSearch => ({
    input: typeof search.input === "string" ? search.input : undefined,
    strict: search.strict === true || search.strict === "true" ? true : undefined,
    validateOnly: search.validateOnly === true || search.validateOnly === "true" ? true : undefined,
  }),
});

const interopRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/interop",
  component: InteropPage,
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

const temporalLabRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/temporal-lab",
  component: TemporalLabPage,
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
