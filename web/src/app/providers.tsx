import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { type ReactNode, useState } from "react";
import { I18nextProvider } from "react-i18next";
import i18n from "@/app/i18n";
import { ThemeProvider } from "@/app/theme";

/** Root providers: TanStack Query + i18next + theme (composition for the app). */
export function Providers({ children }: { children: ReactNode }) {
  const [queryClient] = useState(
    () =>
      new QueryClient({
        defaultOptions: {
          queries: { retry: false, refetchOnWindowFocus: false },
        },
      }),
  );

  return (
    <QueryClientProvider client={queryClient}>
      <I18nextProvider i18n={i18n}>
        <ThemeProvider>{children}</ThemeProvider>
      </I18nextProvider>
    </QueryClientProvider>
  );
}
