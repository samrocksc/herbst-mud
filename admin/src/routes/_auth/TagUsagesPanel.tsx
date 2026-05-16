/* eslint-disable functional/prefer-immutable-types */
import type { Tag, TagUsageReport } from "../../hooks/useTags";
import { Button } from "../../components/Button";

function ColorDot({ color }: { color: string }) {
  const DEFAULT_COLOR = "var(--color-tag-default)";
  const dotStyle = { "--dot-color": color || DEFAULT_COLOR } as React.CSSProperties;
  return (
    <span
      className="inline-block w-3 h-3 rounded-full shrink-0 bg-(--dot-color)"
      style={dotStyle}
    />
  );
}

const safe = <T,>(arr: T[] | undefined): T[] => arr ?? [];

function UsageSection({
  label,
  items,
  badgeClass,
  hrefPrefix,
}: {
  label: string
  items: TagUsageReport["abilities"]
  badgeClass: string
  hrefPrefix: string
}) {
  const list = safe(items);
  if (list.length === 0) return null;
  return (
    <div className="mb-4">
      <h4 className="m-0 mb-2">{label} ({list.length})</h4>
      <ul className="list-none p-0 m-0">
        {list.map((s) => (
          <li key={`${hrefPrefix}-${s.id}`} className="py-1">
            <span className={`badge ${badgeClass} mr-2`}>{s.type}</span>
            <a href={`/${hrefPrefix}?id=${s.id}`} className={`text-${badgeClass.replace("badge-", "")}`}>
              {s.name}
            </a>
          </li>
        ))}
      </ul>
    </div>
  );
}

export function TagUsagesPanel({
  tag,
  report,
  onClose,
}: {
  tag: Tag
  report: TagUsageReport
  onClose: () => void
}) {
  void tag; // tag is available for debugging if needed

  const hasUsages =
    safe(report.abilities).length > 0 || safe(report.factions).length > 0 || safe(report.characters).length > 0;

  return (
    <div className="form-card mt-4">
      <div className="flex justify-between items-center mb-4">
        <h3 className="m-0">
          <ColorDot color={tag.color} />
          <span className="ml-2">{tag.name}</span>
        </h3>
        <Button variant="ghost" size="sm" onClick={onClose} aria-label="Close usages panel">
          ×
        </Button>
      </div>

      {!hasUsages && (
        <div className="empty-state">
          <p className="text-muted">This tag is orphaned — no entities reference it.</p>
        </div>
      )}

      <UsageSection label="Skills" items={report.abilities} badgeClass="badge-accent" hrefPrefix="abilities" />
      <UsageSection label="Factions" items={report.factions} badgeClass="badge-primary" hrefPrefix="factions" />
      <UsageSection label="Characters" items={report.characters} badgeClass="badge-success" hrefPrefix="characters" />
    </div>
  );
}

export function TagUsagesPanelInline({
  report,
}: {
   
  tag: Tag
  report: TagUsageReport
}) {
  const hasUsages =
    safe(report.abilities).length > 0 || safe(report.factions).length > 0 || safe(report.characters).length > 0;

  return (
    <div className="p-3 bg-surface-muted/50 rounded">
      {!hasUsages && <p className="text-muted text-sm m-0">This tag is orphaned — no entities reference it.</p>}
      <UsageSection label="Skills" items={report.abilities} badgeClass="badge-accent" hrefPrefix="abilities" />
      <UsageSection label="Factions" items={report.factions} badgeClass="badge-primary" hrefPrefix="factions" />
      <UsageSection label="Characters" items={report.characters} badgeClass="badge-success" hrefPrefix="characters" />
    </div>
  );
}

export { ColorDot };