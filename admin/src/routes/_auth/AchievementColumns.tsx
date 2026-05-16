import type { Column } from "../../components/DataTable";
import type { Achievement } from "../../hooks/useAchievements";

export const COLUMNS: Column<Achievement>[] = [
  {
    header: "Name",
    accessor: "name",
    render: (_: unknown, row: Achievement) => (
      <strong>{row.icon ? `${row.icon} ` : ""}{row.name}</strong>
    ),
  },
  {
    header: "Description",
    accessor: "description",
  },
  {
    header: "XP",
    accessor: "xp_reward",
    render: (val: unknown) =>
      val ? <span className="talent-effect">{String(val)} XP</span> : <span className="text-muted">—</span>,
  },
  {
    header: "Criteria",
    accessor: "criteria",
    render: (val: unknown) =>
      val ? <span className="text-sm">{String(val)}</span> : <span className="text-muted">—</span>,
  },
  {
    header: "Actions",
    accessor: "_actions",
    render: () => null,
  },
];