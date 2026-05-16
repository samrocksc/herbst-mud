/* eslint-disable no-restricted-syntax */
 
 
 
import { createRouter } from "@tanstack/react-router";
import { routeTree } from "./routeTree.gen";

// Create a new router instance
export const router = createRouter({
  routeTree,
  defaultPreload: "intent",
});

// Augment the module from routeTree.gen.ts instead of redeclaring
type _Router = typeof router;
declare module "./routeTree.gen" {
   
  interface Register {
    router: _Router
  }
}