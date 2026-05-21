import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { StatCard } from "../components/StatCard";
import { StatGrid } from "../components/StatGrid";
import { PageHeader } from "../components/PageHeader";
import { Button } from "../components/Button";
import { useWorlds } from "../hooks/useWorlds";
import { useNPCs } from "../hooks/useNPCs";
import { useRooms } from "../hooks/useRooms";
import { ToolGrid } from "./ToolGrid";
import { useWorldStore } from "../contexts/WorldStoreContext";
import { useEquipmentTemplates } from "../hooks/useEquipmentTemplates";

export const Route = createFileRoute("/dashboard")({
  component: Dashboard,
});

function Dashboard() {
  const navigate = useNavigate();
  const { currentWorld, setWorld } = useWorldStore();

  // Get all worlds to populate the dropdown
  const { data: worlds } = useWorlds();

  // Use world-scoped hooks for counts
  const { rooms } = useRooms();
  const { data: npcsData } = useNPCs();
  const { data: templates, isLoading: templatesLoading } = useEquipmentTemplates();

  // Derived counts
  const roomCount = rooms.length;
  const npcCount = npcsData ? npcsData.length : 0;

  // Handle world switching
  const handleWorldChange = (worldId: string) => {
    setWorld(worldId);
  };

  const handleLogout = () => {
    localStorage.removeItem("token");
    localStorage.removeItem("userId");
    localStorage.removeItem("email");
    localStorage.removeItem("isAdmin");
    navigate({ to: "/login" });
  };

  // Get world name for display
  const currentWorldName = worlds?.find(w => String(w.id) === currentWorld)?.name || currentWorld;

  return (
    <div className="min-h-screen bg-surface text-text p-8">
      <div className="max-w-[1200px] mx-auto">
        <PageHeader
          title="Herbst MUD Admin"
          actions={
            <div className="flex items-center gap-3">
              <select
                value={currentWorld}
                onChange={(e) => handleWorldChange(e.target.value)}
                className="px-3 py-2 bg-surface-muted border border-border rounded text-sm focus:outline-none focus:border-primary"
              >
                {worlds?.map(world => (
                  <option key={world.id} value={world.id}>
                    {world.name}
                  </option>
                ))}
              </select>
              <Button onClick={handleLogout} variant="danger">Logout</Button>
            </div>
          }
        />
        <div className="bg-surface-muted rounded-lg p-6 mb-8">
          <h2 className="m-0 mb-2 text-text">Managing: {currentWorldName}</h2>
          <p className="m-0 text-text-muted">Welcome back! Select a world above to switch contexts.</p>
        </div>
        <StatGrid>
          <StatCard label="Total Rooms" value={roomCount} accent="primary" loading={false} />
          <StatCard label="Active NPCs" value={npcCount} accent="warning" loading={false} />
          <StatCard label="Items" value={templates?.length ?? 0} accent="accent" loading={templatesLoading} />
          <StatCard label="Instances" value={0} accent="primary" loading={false} />
          <StatCard label="Players" value={0} accent="secondary" loading={false} />
          <StatCard label="Skills" value={0} accent="success" loading={false} />
        </StatGrid>
        <h3 className="mb-4 text-text">Admin Tools</h3>
        <ToolGrid />
      </div>
    </div>
  );
}
