import { createFileRoute, Outlet, useLocation } from "@tanstack/react-router";
import { useState } from "react";
import { useTrigger } from "../../hooks/useTriggers";
import { PageHeader } from "../../components/PageHeader";
import { Button } from "../../components/Button";
import { PageContainer } from "../../components/PageContainer";
import { TriggerDetailView } from "./-triggers.$triggerId.detailView";
import { TriggerEditForm } from "./-triggers.$triggerId.editForm";

export const Route = createFileRoute("/_auth/triggers/$triggerId")({
  component: TriggerDetailPage,
});

export function TriggerDetailPage() {
  const triggerId = Route.useParams().triggerId;
  const location = useLocation();
  const { data: trigger, isLoading, error } = useTrigger(Number(triggerId));
  const [editing, setEditing] = useState(false);

  // Render outlet for child routes
  if (location.pathname !== `/triggers/${triggerId}`) {
    return <Outlet />;
  }

  if (isLoading) return <div className="p-8"><PageHeader title="Loading..." backTo="/triggers" /></div>;
  if (error) return <div className="p-8"><PageHeader title="Error" backTo="/triggers" /><div className="text-danger">Failed to load trigger</div></div>;
  if (!trigger) return <div className="p-8"><PageHeader title="Not Found" backTo="/triggers" /><div className="text-danger">Trigger not found</div></div>;

  return (
    <PageContainer>
      <PageHeader
        title={trigger.name}
        backTo="/triggers"
        actions={
          <Button variant={editing ? "secondary" : "primary"} size="sm" onClick={() => setEditing(!editing)}>
            {editing ? "Cancel" : "Edit"}
          </Button>
        }
      />
      {editing ? (
        <TriggerEditForm trigger={trigger} triggerId={Number(triggerId)} onDone={() => setEditing(false)} />
      ) : (
        <TriggerDetailView trigger={trigger} />
      )}
    </PageContainer>
  );
}
