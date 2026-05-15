import { createFileRoute, Outlet } from '@tanstack/react-router';
import { useEffect } from 'react';
import { useNavigate } from '@tanstack/react-router';

export const Route = createFileRoute('/map')({
  component: MapLayout,
});

function MapLayout() {
  const navigate = useNavigate();
  useEffect(() => {
    if (!localStorage.getItem('token')) navigate({ to: '/login' });
  }, [navigate]);

  return <Outlet />;
}
