import type { Tag, TagUsageReport } from '../../hooks/useTags'
import { Button } from '../../components/Button'

function ColorDot({ color }: { color: string }) {
  const DEFAULT_COLOR = 'var(--color-tag-default)'
  const dotStyle = { '--dot-color': color || DEFAULT_COLOR } as React.CSSProperties
  return (
    <span
      className="inline-block w-3 h-3 rounded-full shrink-0 bg-(--dot-color)"
      style={dotStyle}
    />
  )
}

function UsageSection({
  label,
  items,
  badgeClass,
  hrefPrefix,
}: {
  label: string
  items: TagUsageReport['skills']
  badgeClass: string
  hrefPrefix: string
}) {
  if (items.length === 0) return null
  return (
    <div className="mb-4">
      <h4 className="m-0 mb-2">{label} ({items.length})</h4>
      <ul className="list-none p-0 m-0">
        {items.map((s) => (
          <li key={`${hrefPrefix}-${s.id}`} className="py-1">
            <span className={`badge ${badgeClass} mr-2`}>{s.type}</span>
            <a href={`/${hrefPrefix}?id=${s.id}`} className={`text-${badgeClass.replace('badge-', '')}`}>
              {s.name}
            </a>
          </li>
        ))}
      </ul>
    </div>
  )
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
  const hasUsages =
    report.skills.length > 0 || report.factions.length > 0 || report.characters.length > 0

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

      <UsageSection label="Skills" items={report.skills} badgeClass="badge-accent" hrefPrefix="abilities" />
      <UsageSection label="Factions" items={report.factions} badgeClass="badge-primary" hrefPrefix="factions" />
      <UsageSection label="Characters" items={report.characters} badgeClass="badge-success" hrefPrefix="characters" />
    </div>
  )
}

export { ColorDot }