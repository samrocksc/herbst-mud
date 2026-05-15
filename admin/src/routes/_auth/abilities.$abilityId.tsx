import { createFileRoute, useNavigate, Outlet } from '@tanstack/react-router';
import { useState } from 'react';
import { useAbility } from '../../hooks/useAbilities';
import { useLocation } from '@tanstack/react-router';
import { PageHeader } from '../../components/PageHeader';
import { Button } from '../../components/Button';
import { AbilityDetailView } from './-abilities.$abilityId.detailView';
import { AbilityEditForm } from './-abilities.$abilityId.editForm';

export const Route = createFileRoute('/_auth/abilities/$abilityId')({
  component: AbilityDetailPage,
});

function AbilityDetailPage() {
  const abilityId = Route.useParams().abilityId;
  const location = useLocation();
  const navigate = useNavigate();
  const { data: ability, isLoading, error } = useAbility(Number(abilityId));
  const [editing, setEditing] = useState(false);

  // Render outlet for child routes
  if (location.pathname !== `/abilities/${abilityId}`) {
    return <Outlet />;
  }

  if (isLoading) return <div className="p-8"><PageHeader title="Loading..." backTo="/abilities" /></div>;
  if (error) return <div className="p-8"><PageHeader title="Error" backTo="/abilities" /><div className="text-danger">Failed to load ability</div></div>;
  if (!ability) return <div className="p-8"><PageHeader title="Not Found" backTo="/abilities" /><div className="text-danger">Ability not found</div></div>;

  return (
    <div className="p-6 max-w-[1200px] mx-auto">
      <PageHeader
        title={ability.name}
        backTo="/abilities"
        actions={
          <Button variant={editing ? 'secondary' : 'primary'} size="sm" onClick={() => setEditing(!editing)}>
            {editing ? 'Cancel' : 'Edit'}
          </Button>
        }
      />
      {editing ? (
        <AbilityEditForm ability={ability} abilityId={Number(abilityId)} onDone={() => setEditing(false)} />
      ) : (
        <AbilityDetailView ability={ability} />
      )}
    </div>
  );
}
