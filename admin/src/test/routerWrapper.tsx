import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { RouterProvider, createRouter, createRootRoute, createRoute } from "@tanstack/react-router";
import { render } from "@testing-library/react";
import React from "react";

// Create a simple test router
const rootRoute = createRootRoute();
const indexRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/",
  component: () => <div>Root</div>,
});

const routeTree = rootRoute.addChildren([indexRoute]);

const createTestRouter = () => {
  return createRouter({
    routeTree,
  });
};

const testQueryClient = () =>
  new QueryClient({
    defaultOptions: { queries: { retry: false, staleTime: Infinity } },
  });

export function withRouterAndQuery(ui: React.ReactElement) {
  const router = createTestRouter();
  const client = testQueryClient();

  return render(
    <QueryClientProvider client={client}>
      <RouterProvider router={router} />
      {ui}
    </QueryClientProvider>
  );
}