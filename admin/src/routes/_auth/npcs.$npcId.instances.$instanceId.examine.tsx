/* eslint-disable functional/prefer-immutable-types, react-hooks/purity */
import { createFileRoute } from "@tanstack/react-router";
import { useMemo } from "react";
import { useQuery } from "@tanstack/react-query";
import { apiGet } from "../../utils/apiFetch";
import { PageHeader } from "../../components/PageHeader";
import { Section, Field } from "../../components/ExamineFields";

export const Route = createFileRoute("/_auth/npcs/$npcId/instances/$instanceId/examine")({
  component: NPCExaminePage,
});

// ─── Types ───────────────────────────────────────────────────────────────────

type ExamineConfig = Readonly<{
  showName: boolean;
  showDescription: boolean;
  showRace: boolean;
  showLevel: boolean;
  showEquipped: boolean;
  showUnequipped: boolean;
}>;

type EquipmentItem = Readonly<{
  id: number;
  name: string;
  description: string;
  slot: string;
  item_type: string;
}>;

type NPCExamineResponse = Readonly<{
  name: string;
  description: string;
  race: string;
  level: number;
  equipped_items: EquipmentItem[];
  unequipped_items: EquipmentItem[];
}>;

// ─── Default config ──────────────────────────────────────────────────────────

const DEFAULT_CONFIG: ExamineConfig = {
  showName: true,
  showDescription: true,
  showRace: true,
  showLevel: true,
  showEquipped: true,
  showUnequipped: true,
};

// ─── Component ───────────────────────────────────────────────────────────────

function NPCExaminePage() {
  const { npcId, instanceId } = Route.useParams();
  const instanceIdNum = Number(instanceId);

  // Load examine config from game config
  const { data: config } = useExamineConfig();
  const examineConfig: ExamineConfig = config ?? DEFAULT_CONFIG;

  // Load examine details (equipment) via API
  const { data: examineData, isLoading: examineLoading, isError: examineError, error: examineErrorObj } =
    useNPCExamine(instanceIdNum, examineConfig);

  // Combined loading/error state
  const isLoading = examineLoading;
  const isError = examineError;
  const error = examineErrorObj?.message;

  if (isLoading) return <LoadingState />;
  if (isError) return <ErrorState error={error} />;

  // Use examine data for all fields
  const name = examineData?.name || "Unknown";
  const race = examineData?.race || "Unknown";
  const level = examineData?.level ?? 0;
  const description = examineData?.description || "";
  const equippedItems = examineData?.equipped_items ?? [];
  const unequippedItems = examineData?.unequipped_items ?? [];

  return (
    <div className="p-6 max-w-[800px] mx-auto">
      <PageHeader
        title={name}
        showBack
        backTo={`/npcs/${npcId}/instances/${instanceId}`}
      />

      <div className="space-y-6">
        {/* Identity Section */}
        <Section title="Identity">
          {examineConfig.showName && <Field label="Name" value={name} />}
          {examineConfig.showRace && <Field label="Race" value={race} />}
          {examineConfig.showLevel && <Field label="Level" value={String(level)} />}
        </Section>

        {/* Description Section */}
        {examineConfig.showDescription && description && (
          <Section title="Description">
            <p className="text-text-muted text-sm">{description}</p>
          </Section>
        )}

        {/* Equipment Section */}
        {(examineConfig.showEquipped || examineConfig.showUnequipped) && (
          <Section title="Equipment">
            {examineConfig.showEquipped && equippedItems.length > 0 && (
              <div className="mb-4">
                <h3 className="text-sm font-semibold text-text mb-2">Worn Equipment</h3>
                <div className="grid grid-cols-2 gap-2">
                  {equippedItems.map((item) => (
                    <div
                      key={item.id}
                      className="bg-surface-muted p-2 rounded text-sm border border-border"
                    >
                      <div className="font-medium">{item.name}</div>
                      <div className="text-text-muted text-xs">
                        {item.slot} — {item.item_type}
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            )}
            {examineConfig.showUnequipped && unequippedItems.length > 0 && (
              <div>
                <h3 className="text-sm font-semibold text-text mb-2">Unworn Equipment</h3>
                <div className="grid grid-cols-2 gap-2">
                  {unequippedItems.map((item) => (
                    <div
                      key={item.id}
                      className="bg-surface-muted p-2 rounded text-sm border border-border"
                    >
                      <div className="font-medium">{item.name}</div>
                      <div className="text-text-muted text-xs">
                        {item.slot} — {item.item_type}
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            )}
            {examineConfig.showEquipped && equippedItems.length === 0 &&
              examineConfig.showUnequipped && unequippedItems.length === 0 && (
                <div className="text-text-muted text-sm">No equipment</div>
              )}
          </Section>
        )}
      </div>
    </div>
  );
}

// ─── Hooks ───────────────────────────────────────────────────────────────────

function useExamineConfig() {
  return useQuery({
    queryKey: ["examine-config"],
    queryFn: async (): Promise<ExamineConfig> => {
      const data = await apiGet<ExamineConfig>(`${window.location.origin}/api/game-configs/examine_display_config`);
      return data ?? DEFAULT_CONFIG;
    },
  });
}

function useNPCExamine(instanceId: number, config: ExamineConfig) {
  // Build query params based on config
  const params = useMemo(() => {
    const p = new URLSearchParams();
    if (!config.showName) p.append("showName", "false");
    if (!config.showDescription) p.append("showDescription", "false");
    if (!config.showRace) p.append("showRace", "false");
    if (!config.showLevel) p.append("showLevel", "false");
    if (!config.showEquipped) p.append("showEquipped", "false");
    if (!config.showUnequipped) p.append("showUnequipped", "false");
    return p.toString();
  }, [config]);

  return useQuery({
    queryKey: ["npc-examine", instanceId, params],
    queryFn: async (): Promise<NPCExamineResponse | null> => {
      const data = await apiGet<NPCExamineResponse>(
        `${window.location.origin}/api/npc-instances/${instanceId}/examine${params ? `?${params}` : ""}`
      );
      return data ?? null;
    },
  });
}

// ─── UI Components ───────────────────────────────────────────────────────────

function LoadingState() {
  return (
    <div className="p-8">
      <PageHeader title="Loading..." showBack backTo="/npcs" />
      <div className="text-text-muted">Loading NPC details...</div>
    </div>
  );
}

function ErrorState({ error }: { error: string | undefined }) {
  return (
    <div className="p-8">
      <PageHeader title="Error" showBack backTo="/npcs" />
      <div className="text-danger">Failed to load: {error ?? "Unknown error"}</div>
    </div>
  );
}
