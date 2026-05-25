import { createRootRoute, Outlet } from "@tanstack/react-router";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ToastProvider } from "../components/Toast";
import { ErrorBoundary } from "../components/ErrorBoundary";
import { showToast } from "../components/Toast";
import { TopBar } from "../components/TopBar";
import { BottomBar } from "../components/BottomBar";

const queryClient = new QueryClient({
  defaultOptions: {
    mutations: {
      onError: (err) => {
        showToast(err instanceof Error ? err.message : "Operation failed", "error");
      },
    },
  },
});

export const Route = createRootRoute({
  component: RootComponent,
});

function RootComponent() {
  return (
    <QueryClientProvider client={queryClient}>
      <ToastProvider>
        <ErrorBoundary>
          <div className="min-h-[100dvh] bg-surface-muted pt-14 pb-16">
            <TopBar />
            <Outlet />
            <BottomBar />
          </div>
        </ErrorBoundary>
      </ToastProvider>
    </QueryClientProvider>
  );
}
