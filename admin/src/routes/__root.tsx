import { createRootRoute, Outlet } from "@tanstack/react-router";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { useState } from "react";
import { Sidebar } from "../components/Sidebar";
import { MobileNavBar } from "../components/MobileNavBar";
import { MenuIcon } from "../components/icons/MenuIcon";
import { Button } from "../components/Button";
import { ToastProvider } from "../components/Toast";
import { ErrorBoundary } from "../components/ErrorBoundary";
import { showToast } from "../components/Toast";

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
  const [mobileSidebarOpen, setMobileSidebarOpen] = useState(false);
  const [sidebarCollapsed, setSidebarCollapsed] = useState(false);

  return (
    <QueryClientProvider client={queryClient}>
      <ToastProvider>
        <ErrorBoundary>
          <div className="flex min-h-[100dvh] bg-surface-muted">
            {/* Mobile hamburger button */}
            <Button
              variant="ghost"
              size="sm"
              onClick={() => setMobileSidebarOpen(true)}
              aria-label="Open menu"
              className="fixed top-3 left-3 z-40 p-2 bg-surface border border-border text-text-muted hover:bg-surface-muted hover:text-text md:hidden"
            >
              <MenuIcon stroke="currentColor" />
            </Button>

            <Sidebar
              mobileOpen={mobileSidebarOpen}
              onMobileClose={() => setMobileSidebarOpen(false)}
              collapsed={sidebarCollapsed}
              onToggleCollapse={() => setSidebarCollapsed((c) => !c)}
            />

            {/* Mobile backdrop */}
            {mobileSidebarOpen && (
              <div
                className="fixed inset-0 bg-black/30 z-30 md:hidden"
                onClick={() => setMobileSidebarOpen(false)}
              />
            )}

            <div className={`min-h-[100dvh] bg-surface-muted pt-14 pl-14 pr-4 pb-16 md:pt-0 md:pl-0 md:pr-0 md:pb-0 md:ml-16 transition-all duration-300 ease-in-out ${sidebarCollapsed ? 'lg:ml-16' : 'lg:ml-[220px]'}`}>
              <Outlet />
              <MobileNavBar />
            </div>
          </div>
        </ErrorBoundary>
      </ToastProvider>
    </QueryClientProvider>
  );
}
