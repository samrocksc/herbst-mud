import { createRootRoute, Outlet, useLocation } from "@tanstack/react-router";
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
  const location = useLocation();
  const isPublicRoute = location.pathname === "/" || location.pathname === "/login";

  return (
    <QueryClientProvider client={queryClient}>
      <ToastProvider>
        <ErrorBoundary>
          <div className="min-h-[100dvh] bg-surface-muted pt-14 pb-16">
            {!isPublicRoute && <TopBar />}
            <Outlet />
            {!isPublicRoute && <BottomBar />}
          </div>
        </ErrorBoundary>
      </ToastProvider>
    </QueryClientProvider>
  );
}
