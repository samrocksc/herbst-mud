import { useLayoutEffect } from 'react';
import { apiGet } from '../utils/apiFetch';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import { useWorldStore } from '../hooks/useWorldStore';

const API = `${window.location.origin}`;

interface World {
  id: number
  name: string
  title: string
  description: string
  active: boolean
}

function fetchWorldData(worldId: string) {
  if (worldId === 'default') return null;
  return apiGet<{ world?: World; error?: string }>(`${API}/worlds/${worldId}`)
    .then(data => {
      if (!data || data.error) return null;
      return data.world ?? null;
    })
    .catch(() => null);
}

export function WorldTitle() {
  const { currentWorld } = useWorldStore();
  const queryClient = useQueryClient();

  // Refetch world data when world changes
  useLayoutEffect(() => {
    queryClient.invalidateQueries({ queryKey: ['activeWorld', currentWorld] });
  }, [currentWorld, queryClient]);

  // Fetch world data
  const { data: activeWorld, isLoading } = useQuery({
    queryKey: ['activeWorld', currentWorld],
    queryFn: () => fetchWorldData(currentWorld),
    enabled: currentWorld !== 'default',
    staleTime: 0,
  });

  console.log('data', activeWorld, isLoading);

  const title = isLoading ? 'Loading...' : activeWorld?.title ?? 'Herbst MUD';

  return (
    <span className="text-primary font-
      bold text-lg whitespace-nowrap block overflow-hidden text-ellipsis">
      {title}
    </span>
  );
}
