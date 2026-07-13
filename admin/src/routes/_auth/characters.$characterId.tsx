/* eslint-disable functional/prefer-immutable-types, react-hooks/purity */
import { createFileRoute, Link } from "@tanstack/react-router";
import { useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { useCharacter, useUpdateCharacter, useAddCharacterGold, useSpendCharacterGold, type CharacterUpdate } from "../../hooks/useCharacters";
import { apiPost } from "../../utils/apiFetch";
import { PageHeader } from "../../components/PageHeader";
import { EquippedItemsView } from "../../components/EquippedItemsView";
import { ActiveEffectsPanel } from "../../components/ActiveEffectsPanel";
import { ResourceIdField } from "../../components/ResourceIdField";
import { RESOURCE_ENDPOINTS } from "../../utils/resourceEndpoints";
import { AddItemModal } from "./-characters.$characterId.addItemModal";
import { CharacterHistoryPanel } from "../../components/CharacterHistoryPanel";
import { ReclassReraceDialog } from "../../components/ReclassReraceDialog";

function GoldManager({ id, balance }: Readonly<{ id: number; balance: number }>) {
  const [amount, setAmount] = useState(0);
  const addMutation = useAddCharacterGold();
  const spendMutation = useSpendCharacterGold();
  const isPending = addMutation.isPending || spendMutation.isPending;
  const error = (addMutation.error ?? spendMutation.error) as Error | null;

  return (
    <div className="flex items-center gap-2 mt-2">
      <input
        type="number"
        min={0}
        value={amount || ""}
        onChange={(e) => setAmount(parseInt(e.target.value) || 0)}
        className="w-24 p-1 bg-surface border border-border rounded text-text text-sm"
      />
      <button
        onClick={() => addMutation.mutate({ id, amount })}
        disabled={isPending || amount <= 0}
        className="px-2 py-1 bg-success text-white rounded text-xs hover:opacity-90 disabled:opacity-50"
      >
        +
      </button>
      <button
        onClick={() => spendMutation.mutate({ id, amount })}
        disabled={isPending || amount <= 0 || amount > balance}
        className="px-2 py-1 bg-danger text-white rounded text-xs hover:opacity-90 disabled:opacity-50"
      >
        −
      </button>
      {error && <span className="text-danger text-xs">{error.message}</span>}
    </div>
  );
}

export const Route = createFileRoute("/_auth/characters/$characterId")({
  component: CharacterDetail,
});

function CharacterDetail() {
  const { characterId } = Route.useParams();
  const id = Number(characterId);
  const queryClient = useQueryClient();
  const { data: character, isLoading, isError, error } = useCharacter(id);
  const updateMutation = useUpdateCharacter();
  const [editing, setEditing] = useState(false);
  const [showEquipped, setShowEquipped] = useState(true);
  const [addItemOpen, setAddItemOpen] = useState(false);
  const [spawnError, setSpawnError] = useState<string | null>(null);
  const [spawning, setSpawning] = useState(false);
  const [reclassOpen, setReclassOpen] = useState(false);
  const [reraceOpen, setReraceOpen] = useState(false);

  const handleSpawnItem = async (templateId: string) => {
    setSpawning(true);
    setSpawnError(null);
    try {
      await apiPost(`${window.location.origin}/api/item-instances`, {
        equipment_template_id: templateId,
        ownerId: id,
      });
      queryClient.invalidateQueries({ queryKey: ["item-instances"] });
      queryClient.invalidateQueries({ queryKey: ["item-instances", "owner", id] });
      queryClient.invalidateQueries({ queryKey: ["character", id] });
      setAddItemOpen(false);
    } catch (err) {
      setSpawnError(err instanceof Error ? err.message : "Failed to spawn item");
    } finally {
      setSpawning(false);
    }
  };

  if (isLoading) return <div className="p-8 text-text-muted">Loading character...</div>;
  if (isError || !character) return <div className="p-8 text-danger">Error: {error?.message ?? "Not found"}</div>;

  return (
    <div className="p-6 max-w-[800px] mx-auto">
      <PageHeader title={character.name} showBack backTo="/characters" actions={
        <div className="flex gap-2">
          <button
            onClick={() => setReclassOpen(true)}
            className="px-3 py-1.5 bg-warning/20 border border-warning text-warning rounded text-sm hover:bg-warning/30"
            title="Change this character's class (faction)"
          >
            Reclass
          </button>
          <button
            onClick={() => setReraceOpen(true)}
            className="px-3 py-1.5 bg-warning/20 border border-warning text-warning rounded text-sm hover:bg-warning/30"
            title="Change this character's race"
          >
            Rerace
          </button>
          <button
            onClick={() => setAddItemOpen(true)}
            className="px-3 py-1.5 bg-surface border border-border rounded text-sm text-text hover:bg-surface-muted"
          >
            + Add Item
          </button>
          <button
            onClick={() => setShowEquipped(!showEquipped)}
            className="px-3 py-1.5 bg-surface border border-border rounded text-sm text-text hover:bg-surface-muted"
          >
            {showEquipped ? "Hide Equipment" : "Show Equipment"}
          </button>
          <button
            onClick={() => setEditing(!editing)}
            className="px-3 py-1.5 bg-primary text-white rounded text-sm hover:bg-primary-hover"
          >
            {editing ? "Cancel" : "Edit"}
          </button>
          <Link
            to="/characters/$characterId/examine"
            params={{ characterId: String(id) }}
            className="px-3 py-1.5 bg-surface border border-border rounded text-sm text-text hover:bg-surface-muted"
          >
            Examine
          </Link>
        </div>
      } />
      {editing ? (
        <EditForm character={character} onSave={(update) => {
          updateMutation.mutate({ id, update }, { onSuccess: () => setEditing(false) });
        }} />
      ) : (
        <DetailView character={character} />
      )}
      {showEquipped && <EquippedItemsView characterId={character.id} characterRace={character.race} />}
      <ActiveEffectsPanel characterId={character.id} />

      {/* Class & Race History */}
      <div className="mt-6">
        <h3 className="text-sm font-semibold text-text mb-2 pb-1 border-b border-border">History</h3>
        <CharacterHistoryPanel characterId={character.id} />
      </div>

      <AddItemModal open={addItemOpen} onClose={() => { setAddItemOpen(false); setSpawnError(null); }}
        onSpawn={handleSpawnItem} isLoading={spawning} error={spawnError} />

      {/* Reclass / Rerace dialogs */}
      <ReclassReraceDialog
        open={reclassOpen}
        mode="reclass"
        characterId={character.id}
        currentName={character.name}
        currentClassName={character.class}
        currentRaceName={character.race}
        onClose={() => setReclassOpen(false)}
      />
      <ReclassReraceDialog
        open={reraceOpen}
        mode="rerace"
        characterId={character.id}
        currentName={character.name}
        currentClassName={character.class}
        currentRaceName={character.race}
        onClose={() => setReraceOpen(false)}
      />
    </div>
  );
}

 
function DetailView({ character }: { character: NonNullable<ReturnType<typeof useCharacter>["data"]> }) {
  const lastSeen = character.lastSeenAt ? new Date(character.lastSeenAt).toLocaleString() : "Never";
   
  const isOnline = character.lastSeenAt && new Date(character.lastSeenAt) > new Date(Date.now() - 15 * 60 * 1000);

  return (
    <div className="space-y-6">
      <Section title="Identity">
        <Field label="ID" value={String(character.id)} />
        <Field label="Name" value={character.name} />
        <Field label="Race" value={character.race} />
        <Field label="Class" value={character.class} />
        <Field label="Gender" value={character.gender || "—"} />
        <Field label="NPC" value={character.isNPC ? "Yes" : "No"} />
        <Field label="Admin" value={character.is_admin ? "Yes" : "No"} />
        <Field label="Test" value={character.is_test ? "Yes" : "No"} />
        <Field label="Immortal" value={character.is_immortal ? "Yes" : "No"} />
      </Section>

      <Section title="Location">
        <Field label="Current Room" value={`#${character.currentRoomId}`} />
        <Field label="Starting Room" value={`#${character.startingRoomId}`} />
        <Field label="Respawn Room" value={`#${character.respawnRoomId}`} />
        <Field label="World" value={character.currentWorld} />
        <Field label="Last Seen" value={lastSeen} />
        <Field label="Status" value={isOnline ? "Online" : "Offline"} />
      </Section>

      <Section title="Vitals">
        <Field label="HP" value={`${character.hitpoints} / ${character.max_hitpoints}`} />
        <Field label="Stamina" value={`${character.stamina} / ${character.max_stamina}`} />
        <Field label="Mana" value={`${character.mana} / ${character.max_mana}`} />
      </Section>

      <Section title="Progression">
        <Field label="Level" value={String(character.level)} />
        <Field label="XP" value={String(character.xp)} />
      </Section>

      <Section title="Currency">
        <Field label="Gold Credits" value={String(character.gold_credits ?? 0)} />
        <GoldManager id={character.id} balance={character.gold_credits ?? 0} />
      </Section>

      <Section title="Stats">
        <Field label="STR" value={String(character.strength)} />
        <Field label="DEX" value={String(character.dexterity)} />
        <Field label="CON" value={String(character.constitution)} />
        <Field label="INT" value={String(character.intelligence)} />
        <Field label="WIS" value={String(character.wisdom)} />
      </Section>

      {character.description && (
        <Section title="Description">
          <p className="text-text-muted text-sm">{character.description}</p>
        </Section>
      )}
    </div>
  );
}

function EditForm({ character, onSave }: { character: NonNullable<ReturnType<typeof useCharacter>["data"]>, onSave: (update: CharacterUpdate) => void }) {
  const [form, setForm] = useState<CharacterUpdate>({
    name: character.name,
    currentRoomId: character.currentRoomId,
    startingRoomId: character.startingRoomId,
    respawnRoomId: character.respawnRoomId,
    level: character.level,
    xp: character.xp,
    gold_credits: character.gold_credits ?? 0,
    hitpoints: character.hitpoints,
    maxHitpoints: character.max_hitpoints,
    stamina: character.stamina,
    maxStamina: character.max_stamina,
    mana: character.mana,
    maxMana: character.max_mana,
    gender: character.gender,
    description: character.description,
    isNPC: character.isNPC,
    isAdmin: character.is_admin,
    isTest: character.is_test,
  });

  const numField = (key: keyof CharacterUpdate, label: string) => (
    <div className="flex items-center gap-2 mb-2">
      <label className="text-text-muted text-sm w-28 shrink-0">{label}</label>
      <input
        type="number"
        value={form[key] as number ?? 0}
        onChange={(e) => setForm({ ...form, [key]: parseInt(e.target.value) || 0 })}
        className="w-24 p-1 bg-surface border border-border rounded text-text text-sm"
      />
    </div>
  );

  return (
    <div className="space-y-6">
      <Section title="Identity">
        <div className="flex items-center gap-2 mb-2">
          <label className="text-text-muted text-sm w-28 shrink-0">Name</label>
          <input
            type="text"
            value={form.name ?? ""}
            onChange={(e) => setForm({ ...form, name: e.target.value })}
            className="w-48 p-1 bg-surface border border-border rounded text-text text-sm"
          />
        </div>
        <div className="flex items-center gap-2 mb-2">
          <label className="text-text-muted text-sm w-28 shrink-0">Gender</label>
          <input
            type="text"
            value={form.gender ?? ""}
            onChange={(e) => setForm({ ...form, gender: e.target.value })}
            className="w-32 p-1 bg-surface border border-border rounded text-text text-sm"
          />
        </div>
        <div className="flex items-center gap-2 mb-2">
          <label className="text-text-muted text-sm w-28 shrink-0">Description</label>
          <input
            type="text"
            value={form.description ?? ""}
            onChange={(e) => setForm({ ...form, description: e.target.value })}
            className="w-full p-1 bg-surface border border-border rounded text-text text-sm"
          />
        </div>
        <div className="flex items-center gap-4 mb-2">
          <label className="flex items-center gap-1 text-sm text-text-muted cursor-pointer">
            <input type="checkbox" checked={form.isNPC ?? false} onChange={(e) => setForm({ ...form, isNPC: e.target.checked })} className="accent-primary" />
            NPC
          </label>
          <label className="flex items-center gap-1 text-sm text-text-muted cursor-pointer">
            <input type="checkbox" checked={form.isAdmin ?? false} onChange={(e) => setForm({ ...form, isAdmin: e.target.checked })} className="accent-primary" />
            Admin
          </label>
          <label className="flex items-center gap-1 text-sm text-text-muted cursor-pointer">
            <input type="checkbox" checked={form.isTest ?? false} onChange={(e) => setForm({ ...form, isTest: e.target.checked })} className="accent-primary" />
            Test
          </label>
        </div>
      </Section>

<Section title="Location">
        <ResourceIdField
          label="Current Room"
          value={form.currentRoomId}
          onChange={(id) => setForm({ ...form, currentRoomId: Number(id) })}
          {...RESOURCE_ENDPOINTS.rooms}
        />
        <div className="mt-2">
          <ResourceIdField
            label="Starting Room"
            value={form.startingRoomId}
            onChange={(id) => setForm({ ...form, startingRoomId: Number(id) })}
            {...RESOURCE_ENDPOINTS.rooms}
          />
        </div>
        <div className="mt-2">
          <ResourceIdField
            label="Respawn Room"
            value={form.respawnRoomId}
            onChange={(id) => setForm({ ...form, respawnRoomId: Number(id) })}
            {...RESOURCE_ENDPOINTS.rooms}
          />
        </div>
      </Section>

      <Section title="Vitals">
        {numField("hitpoints", "HP")}
        {numField("maxHitpoints", "Max HP")}
        {numField("stamina", "Stamina")}
        {numField("maxStamina", "Max Stamina")}
        {numField("mana", "Mana")}
        {numField("maxMana", "Max Mana")}
      </Section>

      <Section title="Progression">
        {numField("level", "Level")}
        {numField("xp", "XP")}
        {numField("gold_credits", "Gold Credits")}
      </Section>

      <button
        onClick={() => onSave(form)}
        className="px-4 py-2 bg-accent text-white rounded text-sm hover:bg-accent-hover"
      >
        Save Changes
      </button>
    </div>
  );
}

function Section({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <div>
      <h3 className="text-sm font-semibold text-text mb-2 pb-1 border-b border-border">{title}</h3>
      <div className="space-y-0.5">{children}</div>
    </div>
  );
}

function Field({ label, value }: { label: string; value: string }) {
  return (
    <div className="flex gap-2 text-sm">
      <span className="text-text-muted w-28 shrink-0">{label}</span>
      <span className="text-text">{value}</span>
    </div>
  );
}