 
import { useLayoutEffect } from "react";
import { apiGet } from "../utils/apiFetch";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { useWorldStore } from "../contexts/WorldStoreContext";
import type { World } from "../hooks/useWorlds";

const API = `${window.location.origin}`;

const fetchWorldData = async (worldId: string): Promise<World | null> => {
  if (worldId === "default") return null;
  try {
    const data = await apiGet<{ world?: World; error?: string }>(`${API}/worlds/${worldId}`);
    if (!data || data.error) return null;
    return data.world ?? null;
  } catch {
    return null;
  }
};

export function WorldTitle() {
  const { currentWorld } = useWorldStore();
  const queryClient = useQueryClient();

  useLayoutEffect(() => {
    queryClient.invalidateQueries({ queryKey: ["activeWorld", currentWorld] });
  }, [currentWorld, queryClient]);

  const { data: activeWorld, isLoading } = useQuery({
    queryKey: ["activeWorld", currentWorld],
    queryFn: () => fetchWorldData(currentWorld),
    enabled: currentWorld !== "default",
    staleTime: 0,
  });

  const title = isLoading ? "Loading..." : activeWorld?.title ?? "Herbst MUD";

  return (
    <span className="text-primary font-bold text-lg whitespace-nowrap block overflow-hidden text-ellipsis">
      {title}
    </span>
  );
}
