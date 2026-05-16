/* eslint-disable functional/prefer-immutable-types, functional/no-mixed-types */
import { Fragment, type ReactNode } from "react";

// ─── Column definition ──────────────────────────────────────────────────────

type Column<T> = Readonly<{
  header: string
  accessor: string
  render?: (value: unknown, row: Readonly<T>) => ReactNode
  className?: string
  align?: "left" | "center" | "right"
}>

// ─── Component props ────────────────────────────────────────────────────────

type DataTableProps<T> = Readonly<{
  columns: Column<T>[]
  data: ReadonlyArray<T>
  getKey: (row: Readonly<T>) => string | number
  onRowClick?: (row: Readonly<T>) => void
  expandedRow?: (row: Readonly<T>) => ReactNode
  className?: string
  emptyMessage?: string
  variant?: "default" | "dark"
}>

// ─── Helpers ─────────────────────────────────────────────────────────────────

/** Read any leaf value from an object via a dot-notation path. */
function getValue<T>(row: T, accessor: string): unknown {
  return accessor.split(".").reduce<unknown>((acc, key) => {
    if (acc == null || typeof acc !== "object") return undefined;
    return (acc as Record<string, unknown>)[key];
  }, row);
}

// Default cell renderer: '-' for nullish, green/red badge for booleans
function DefaultCell({ value }: { value: unknown }): ReactNode {
  if (value === null || value === undefined || value === "") return <span className="text-muted">—</span>;
  if (typeof value === "boolean") {
    return value
      ? <span className="badge badge-success">Yes</span>
      : <span className="badge badge-neutral">No</span>;
  }
  return <>{String(value)}</>;
}

// ─── Component ───────────────────────────────────────────────────────────────

export function DataTable<T>({
  columns,
  data,
  getKey,
  onRowClick,
  expandedRow,
  className = "",
  emptyMessage = "No records found.",
  variant = "default",
}: DataTableProps<T>) {
  const tableClass = variant === "dark" ? "table table-dark" : "table";

  const alignClass = (align?: "left" | "center" | "right") =>
    align === "center" ? "text-center" : align === "right" ? "text-right" : "text-left";

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
            data.map((row: Readonly<T>) => {
              const key = getKey(row);
              const expanded = expandedRow?.(row);
              return (
                <Fragment key={key}>
                  <tr onClick={onRowClick ? () => onRowClick(row) : undefined}
                      className={onRowClick ? "clickable-row" : undefined}>
                    {columns.map((col: Column<T>) => {
                      const raw = getValue(row, col.accessor);
                      return (
                        <td
                          key={col.accessor}
                          className={[col.className, alignClass(col.align)].filter(Boolean).join(" ") || undefined}
                        >
                          {col.render ? col.render(raw, row) : <DefaultCell value={raw} />}
                        </td>
                      );
                    })}
                  </tr>
                  {expanded && (
                    <tr className="expanded-row">
                      <td colSpan={columns.length}>{expanded}</td>
                    </tr>
                  )}
                </Fragment>
              );
            })
          )}
        </tbody>
      </table>
    </div>
  );
}

export type { Column };