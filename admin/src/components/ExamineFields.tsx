/* eslint-disable functional/prefer-immutable-types */
import React from "react";

// ─── Section ─────────────────────────────────────────────────────────────────

export function Section({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <div>
      <h3 className="text-sm font-semibold text-text mb-2 pb-1 border-b border-border">{title}</h3>
      <div className="space-y-0.5">{children}</div>
    </div>
  );
}

// ─── Field ───────────────────────────────────────────────────────────────────

export function Field({ label, value }: { label: string; value: string }) {
  return (
    <div className="flex gap-2 text-sm">
      <span className="text-text-muted w-28 shrink-0">{label}</span>
      <span className="text-text">{value}</span>
    </div>
  );
}

// ─── Equipment Item Display ─────────────────────────────────────────────────

export function EquipmentItemDisplay({ item }: { item: Readonly<{ id: number; name: string; slot: string; item_type: string }> }) {
  return (
    <div className="bg-surface-muted p-2 rounded text-sm border border-border">
      <div className="font-medium">{item.name}</div>
      <div className="text-text-muted text-xs">
        {item.slot} — {item.item_type}
      </div>
    </div>
  );
}
