import type { ReactNode } from 'react'

// ─── Column definition ──────────────────────────────────────────────────────

type Column<T> = Readonly<{
  /** Static header label */
  header: string
  /** Dot-notation path into T, e.g. 'name' or 'stats.xp' */
  accessor: string
  /** Override the default cell render. Receives (value, row). */
  render?: (value: unknown, row: T) => ReactNode
  /** Additional CSS class for this <td> */
  className?: string
  /** Text alignment for this column */
  align?: 'left' | 'center' | 'right'
}>

// ─── Component props ────────────────────────────────────────────────────────

type DataTableProps<T> = Readonly<{
  /** Column definitions */
  columns: Column<T>[]
  /** Row data */
  data: T[]
  /** Unique key extractor */
  getKey: (row: T) => string | number
  /** Optional row-level click handler (enables clickable-row hover state) */
  onRowClick?: (row: T) => void
  /** Extra CSS class on the wrapper */
  className?: string
  /** Override the empty-state message */
  emptyMessage?: string
  /** Visual style variant */
  variant?: 'default' | 'dark'
}>

// ─── Helpers ─────────────────────────────────────────────────────────────────

/** Read any leaf value from an object via a dot-notation path. */
function getValue<T>(row: T, accessor: string): unknown {
  return accessor.split('.').reduce<unknown>((acc, key) => {
    if (acc == null || typeof acc !== 'object') return undefined
    return (acc as Record<string, unknown>)[key]
  }, row)
}

// Default cell renderer: '-' for nullish, green/red badge for booleans
function DefaultCell({ value }: { value: unknown }): ReactNode {
  if (value === null || value === undefined || value === '') return <span className="text-muted">—</span>
  if (typeof value === 'boolean') {
    return value
      ? <span className="badge badge-success">Yes</span>
      : <span className="badge badge-neutral">No</span>
  }
  return <>{String(value)}</>
}

// ─── Component ───────────────────────────────────────────────────────────────

export function DataTable<T>({
  columns,
  data,
  getKey,
  onRowClick,
  className = '',
  emptyMessage = 'No records found.',
  variant = 'default',
}: DataTableProps<T>) {
  const tableClass = variant === 'dark' ? 'table table-dark' : 'table'

  const alignClass = (align?: 'left' | 'center' | 'right') =>
    align === 'center' ? 'text-center' : align === 'right' ? 'text-right' : 'text-left'

  return (
    <div className={`table-container ${className}`}>
      <table className={tableClass}>
        <thead>
          <tr>
            {columns.map((col) => (
              <th
                key={col.accessor}
                className={alignClass(col.align)}
              >
                {col.header}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {data.length === 0 ? (
            <tr>
              <td colSpan={columns.length}>
                <div className="empty-state">
                  <p>{emptyMessage}</p>
                </div>
              </td>
            </tr>
          ) : (
            data.map((row: T) => {
              const key = getKey(row)
              return (
                <tr key={key} onClick={onRowClick ? () => onRowClick(row) : undefined}
                    className={onRowClick ? 'clickable-row' : undefined}>
                  {columns.map((col: Column<T>) => {
                    const raw = getValue(row, col.accessor)
                    return (
                      <td
                        key={col.accessor}
                        className={[col.className, alignClass(col.align)].filter(Boolean).join(' ') || undefined}
                      >
                        {col.render ? col.render(raw, row) : <DefaultCell value={raw} />}
                      </td>
                    )
                  })}
                </tr>
              )
            })
          )}
        </tbody>
      </table>
    </div>
  )
}

export type { Column }