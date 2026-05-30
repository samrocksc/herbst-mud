import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { render } from "@testing-library/react";

const testQueryClient = () =>
  new QueryClient({
    defaultOptions: { queries: { retry: false, staleTime: Infinity } },
  });

export function withQuery(ui: React.ReactElement) {
  const client = testQueryClient();
  return {
    client,
    ...render(<QueryClientProvider client={client}>{ui}</QueryClientProvider>),
  };
}

export function QueryClientWrapper({ children }: { children: React.ReactNode }) {
  const client = testQueryClient();
  return <QueryClientProvider client={client}>{children}</QueryClientProvider>;
}
