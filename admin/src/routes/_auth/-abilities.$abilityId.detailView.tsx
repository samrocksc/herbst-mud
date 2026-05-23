import type { Ability } from "../../hooks/useAbilities";

export function AbilityDetailView({ ability }: Readonly<{ ability: Ability }>) {
  return (
    <div className="bg-surface-muted rounded-lg p-6 border border-border">
      <h2 className="mt-0 mb-4 text-text text-lg font-semibold">Ability Details</h2>
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
        <DetailField label="ID" value={String(ability.id)} />
        <DetailField label="Name" value={ability.name} />
        <DetailField label="Type" value={ability.ability_type} />
        <DetailField label="Class" value={ability.ability_class} />
        <DetailField label="Slug" value={ability.slug} />
        <DetailField label="Required Tag" value={ability.required_tag || "—"} />
        <DetailField label="Level Req" value={ability.requirements || "—"} />
        <DetailField label="Cost" value={String(ability.cost)} />
        <DetailField label="Cooldown" value={ability.cooldown_seconds > 0 ? `${ability.cooldown_seconds}s` : "—"} />
        <DetailField label="Mana Cost" value={ability.mana_cost > 0 ? String(ability.mana_cost) : "—"} />
        <DetailField label="Stamina Cost" value={ability.stamina_cost > 0 ? String(ability.stamina_cost) : "—"} />
        <DetailField label="HP Cost" value={ability.hp_cost > 0 ? String(ability.hp_cost) : "—"} />
        {ability.proc_chance > 0 && (
          <DetailField label="Proc Chance" value={String(ability.proc_chance)} />
        )}
        {ability.proc_event && (
          <DetailField label="Proc Event" value={ability.proc_event} />
        )}
        {ability.faction_skills !== null && (
          <DetailField label="Faction Skills" value={String(ability.faction_skills)} />
        )}
      </div>
      {ability.description && (
        <p className="text-text text-sm mt-4">{ability.description}</p>
      )}
    </div>
  );
}

function DetailField({ label, value }: Readonly<{ label: string; value: string }>) {
  return (
    <div>
      <span className="text-text-muted text-xs block mb-0.5">{label}</span>
      <span className="text-text text-sm font-medium">{value}</span>
    </div>
  );
}
