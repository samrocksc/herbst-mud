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

// ─── Cell renderer ───────────────────────────────────────────────────────────

function renderCell<T>(col: Column<T>, row: T): ReactNode {
  const raw = getValue(row, col.accessor);
  if (col.render) return col.render(raw, row);
  return <DefaultCell value={raw} />;
}

// ─── Mobile card view ────────────────────────────────────────────────────────

function MobileCardList<T>({
  columns,
  data,
  getKey,
  onRowClick,
  expandedRow,
}: Omit<DataTableProps<T>, "className" | "emptyMessage" | "variant">) {
  const metaCols = columns.filter((c) => c.header !== "");
  const actionCols = columns.filter((c) => c.header === "");

  return (
    <div className="space-y-3">
      {data.map((row) => {
        const key = getKey(row);
        const clickable = !!onRowClick;
        return (
          <Fragment key={key}>
            <div
              className={[
                "bg-surface border border-border rounded-lg p-3",
                clickable ? "cursor-pointer hover:bg-surface-muted" : "",
                "transition-colors",
              ].join(" ")}
              onClick={clickable ? () => onRowClick!(row) : undefined}
            >
              {/* Meta fields */}
              {metaCols.map((col) => (
                <div key={col.accessor} className="flex items-start justify-between py-1 gap-3">
                  <span className="text-text-muted text-xs shrink-0">{col.header}</span>
                  <div className={[`text-sm text-right min-w-0`, col.align === "center" ? "text-center" : col.align === "right" ? "text-right" : "text-left"].join(" ")}>
                    {renderCell(col, row)}
                  </div>
                </div>
              ))}

              {/* Action buttons */}
              {actionCols.length > 0 && (
                <div className="flex justify-end gap-2 pt-2 mt-2 border-t border-border">
                  {actionCols.map((col) => (
                    <div key={col.accessor}>{renderCell(col, row)}</div>
                  ))}
                </div>
              )}
            </div>

            {/* Expanded row */}
            {expandedRow && <div className="px-1">{expandedRow(row)}</div>}
          </Fragment>
        );
      })}
    </div>
  );
}

// ─── Desktop table view ──────────────────────────────────────────────────────

function DesktopTable<T>({
  columns,
  data,
  getKey,
  onRowClick,
  expandedRow,
  variant,
}: Omit<DataTableProps<T>, "className" | "emptyMessage">) {
  const tableClass = variant === "dark" ? "table table-dark" : "table";

  const alignClass = (align?: "left" | "center" | "right") =>
    align === "center" ? "text-center" : align === "right" ? "text-right" : "text-left";

  return (
    <div className="table-container">
      <table className={tableClass}>
        <thead>
          <tr>
            {columns.map((col) => (
              <th key={col.accessor} className={alignClass(col.align)}>
                {col.header}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {data.map((row) => {
            const key = getKey(row);
            const clickable = !!onRowClick;
            return (
              <Fragment key={key}>
                <tr
                  className={clickable ? "clickable-row" : ""}
                  onClick={clickable ? () => onRowClick!(row) : undefined}
                >
                  {columns.map((col) => (
                    <td
                      key={col.accessor}
                      className={`${alignClass(col.align)} ${col.className ?? ""}`}
                    >
                      {renderCell(col, row)}
                    </td>
                  ))}
                </tr>
                {expandedRow && (
                  <tr className="expanded-row">
                    <td colSpan={columns.length}>{expandedRow(row)}</td>
                  </tr>
                )}
              </Fragment>
            );
          })}
        </tbody>
      </table>
    </div>
  );
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
  if (data.length === 0) {
    return (
      <div className={`p-6 sm:p-8 text-center text-text-muted text-sm ${className}`}>
        {emptyMessage}
      </div>
    );
  }

  return (
    <div className={className}>
      {/* Mobile card list */}
      <div className="sm:hidden">
        <MobileCardList
          columns={columns}
          data={data}
          getKey={getKey}
          onRowClick={onRowClick}
          expandedRow={expandedRow}
        />
      </div>

      {/* Desktop table */}
      <div className="hidden sm:block">
        <DesktopTable
          columns={columns}
          data={data}
          getKey={getKey}
          onRowClick={onRowClick}
          expandedRow={expandedRow}
          variant={variant}
        />
      </div>
    </div>
  );
}

export type { Column };
