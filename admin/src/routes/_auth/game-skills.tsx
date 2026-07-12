/* eslint-disable functional/immutable-data */
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import {
  useGameSkills,
  useDeleteGameSkill,
  type GameSkill,
} from "../../hooks/useGameSkills";
import { PageHeader } from "../../components/PageHeader";
import { DataTable, type Column } from "../../components/DataTable";
import { Button } from "../../components/Button";
import { PageContainer } from "../../components/PageContainer";
import { FilterBar } from "../../components/FilterBar";
import { DeleteConfirmation } from "../../components/DeleteConfirmation";
import { GameSkillForm } from "./-GameSkillForm";

export const Route = createFileRoute("/_auth/game-skills")({
  component: GameSkillsManagement,
});

const CATEGORIES = ["weapon", "armor", "craft", "magic"];

const COLUMNS: Column<GameSkill>[] = [
  {
    header: "Name",
    accessor: "name",
    render: (_, row) => (
      <span className="font-bold text-primary">{row.name}</span>
    ),
  },
  {
    header: "Display Name",
    accessor: "display_name",
    render: (val: unknown) => (
      <span className="text-text">{String(val ?? "")}</span>
    ),
  },
  {
    header: "Description",
    accessor: "description",
    render: (val: unknown) => (
      <span className="text-text-muted text-xs" title={String(val ?? "")}>
        {String(val ?? "").slice(0, 80)}
        {String(val ?? "").length > 80 ? "…" : ""}
      </span>
    ),
  },
  {
    header: "Category",
    accessor: "category",
    render: (val: unknown) => (
      <span className="text-xs px-2 py-0.5 rounded bg-primary/15 text-text border border-primary/30 capitalize">
        {String(val)}
      </span>
    ),
  },
  {
    header: "Max Level",
    accessor: "max_level",
    render: (val: unknown) => (
      <span className="text-text">{String(val ?? "")}</span>
    ),
  },
  {
    header: "XP Curve",
    accessor: "xp_curve_mode",
    render: (val: unknown) => (
      <span className="text-xs px-2 py-0.5 rounded bg-accent/15 text-text border border-accent/30 capitalize">
        {String(val)}
      </span>
    ),
  },
];

export function GameSkillsManagement() {
  const [filterCategory, setFilterCategory] = useState<string>("");
  const [showForm, setShowForm] = useState(false);
  const [editingSkill, setEditingSkill] = useState<GameSkill | null>(null);
  const [deletingSkill, setDeletingSkill] = useState<GameSkill | null>(null);
  const navigate = useNavigate();

  const { data: skills, isLoading, error } = useGameSkills({
    category: filterCategory || undefined,
  });
  const deleteMutation = useDeleteGameSkill();

  const handleEdit = (skill: GameSkill) => {
    setEditingSkill(skill);
    setShowForm(true);
  };

  const handleAdd = () => {
    setEditingSkill(null);
    setShowForm(true);
  };

  const handleDelete = async () => {
    if (!deletingSkill) return;
    try {
      await deleteMutation.mutateAsync(deletingSkill.id);
      setDeletingSkill(null);
    } catch {
      /* error is in mutation state */
    }
  };

  const columns: Column<GameSkill>[] = [
    ...COLUMNS,
    {
      header: "Actions",
      accessor: "_actions",
      render: (_: unknown, row: GameSkill) => (
        <div className="flex gap-2 justify-end">
          <Button variant="accent" size="sm" onClick={() => handleEdit(row)}>
            Edit
          </Button>
          <Button variant="danger" size="sm" onClick={() => setDeletingSkill(row)}>
            Delete
          </Button>
        </div>
      ),
    },
  ];

  if (showForm) {
    return (
      <PageContainer>
        <PageHeader
          title={editingSkill ? "Edit Skill" : "Create Skill"}
          showBack
          backTo="/game-skills"
          backLabel="← Game Skills"
        />
        <GameSkillForm
          skill={editingSkill}
          onSubmit={() => {
            setShowForm(false);
            setEditingSkill(null);
            navigate({ to: "/game-skills" });
          }}
          onCancel={() => {
            setShowForm(false);
            setEditingSkill(null);
          }}
        />
      </PageContainer>
    );
  }

  if (isLoading) return <div className="loading">Loading game skills...</div>;
  if (error) return <div className="error">Failed to load game skills: {error.message}</div>;

  return (
    <PageContainer>
      <PageHeader
        title="Game Skills"
        backTo="/dashboard"
        actions={
          <Button variant="primary" onClick={handleAdd}>
            + Add Skill
          </Button>
        }
      />
      <p className="text-sm text-muted mb-4">
        Game skills are leveled proficiencies — Blades, Heavy Armor, Crafting, etc. — that gate
        equipment use and provide bonuses as they level up. Each skill has a category, max level,
        and an XP curve that defines how XP maps to skill levels.
      </p>

      <FilterBar
        showClear={!!filterCategory}
        onClear={() => setFilterCategory("")}
      >
        <div className="flex flex-col gap-1">
          <label className="text-xs text-text-muted">Category:</label>
          <select
            value={filterCategory}
            onChange={(e) => setFilterCategory(e.target.value)}
            className="px-3 py-2 bg-surface border border-border rounded text-sm text-text focus:outline-none focus:border-primary"
          >
            <option value="">All Categories</option>
            {CATEGORIES.map((cat) => (
              <option key={cat} value={cat} className="capitalize">
                {cat.charAt(0).toUpperCase() + cat.slice(1)}
              </option>
            ))}
          </select>
        </div>
      </FilterBar>

      <DataTable
        columns={columns}
        data={skills ?? []}
        getKey={(row) => row.id}
        emptyMessage={
          filterCategory
            ? "No skills match this category filter."
            : "No game skills found. Create your first skill!"
        }
      />

      <DeleteConfirmation
        open={!!deletingSkill}
        title="Delete Skill"
        message={
          deletingSkill
            ? `Are you sure you want to delete "${deletingSkill.display_name || deletingSkill.name}"? This action cannot be undone.`
            : ""
        }
        onConfirm={handleDelete}
        onCancel={() => setDeletingSkill(null)}
        isLoading={deleteMutation.isPending}
      />
    </PageContainer>
  );
}