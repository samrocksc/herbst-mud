import { createFileRoute, Link } from "@tanstack/react-router";
import { useState } from "react";
import { PageHeader } from "../../components/PageHeader";
import { DataTable, type Column } from "../../components/DataTable";
import { Button } from "../../components/Button";
import { useRecipes, useDeleteRecipe, type Recipe } from "../../hooks/useRecipes";
import { useWorldStore } from "../../contexts/WorldStoreContext";

export const Route = createFileRoute("/_auth/recipes")({
  component: RecipesManagement,
});

function DeleteConfirmation({
  recipe,
  onConfirm,
  onCancel,
  isLoading,
}: Readonly<{
  recipe: Recipe
  onConfirm: () => void
  onCancel: () => void
  isLoading: boolean
}>) {
  return (
    <div className="modal-overlay" onClick={onCancel}>
      <div className="modal-content modal-sm" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <h3>Delete Recipe</h3>
          <Button variant="ghost" size="sm" onClick={onCancel} aria-label="Close">×</Button>
        </div>
        <div className="modal-body">
          <p>Are you sure you want to delete <strong>{recipe.display_name}</strong>?</p>
          <p className="text-muted">This action cannot be undone.</p>
        </div>
        <div className="modal-footer">
          <Button variant="danger" onClick={onConfirm} disabled={isLoading}>
            {isLoading ? "Deleting..." : "Delete"}
          </Button>
          <Button variant="secondary" onClick={onCancel}>Cancel</Button>
        </div>
      </div>
    </div>
  );
}

function RecipesManagement() {
  const [deletingRecipe, setDeletingRecipe] = useState<Recipe | null>(null);

  const { currentWorld } = useWorldStore();
  const { data: recipes, isLoading, error } = useRecipes({ world_id: currentWorld });
  const deleteMutation = useDeleteRecipe();

  const handleDelete = async () => {
    if (!deletingRecipe) return;
    try {
      await deleteMutation.mutateAsync(deletingRecipe.name);
      setDeletingRecipe(null);
    } catch { /* error is in mutation state */ }
  };

  const columns: Column<Recipe>[] = [
    {
      header: "Name",
      accessor: "name",
      render: (_: unknown, row: Recipe) => (
        <Link
          to="/recipes/$recipeName"
          params={{ recipeName: row.name }}
          className="text-primary no-underline hover:underline font-bold"
        >
          {row.name}
        </Link>
      ),
    },
    { header: "Display Name", accessor: "display_name" },
    { header: "Station", accessor: "required_station_tag" },
    { header: "Class", accessor: "required_class", render: (val: unknown) => String(val) || "Any" },
    {
      header: "Inputs → Outputs",
      accessor: "_inputs_outputs",
      render: (_: unknown, row: Recipe) => {
        const inCount = row.inputs?.length ?? 0;
        const outCount = row.outputs?.length ?? 0;
        return `${inCount} → ${outCount}`;
      },
    },
    { header: "Time (s)", accessor: "craft_time_secs" },
    {
      header: "Actions",
      accessor: "_actions",
      render: (_: unknown, row: Recipe) => (
        <>
          <Link
            to="/recipes/$recipeName/edit"
            params={{ recipeName: row.name }}
            className="text-accent no-underline hover:underline text-sm mr-2"
          >
            Edit
          </Link>
          <span
            className="text-danger cursor-pointer text-sm hover:underline"
            onClick={() => setDeletingRecipe(row)}
          >
            Delete
          </span>
        </>
      ),
    },
  ];

  if (isLoading) return <div className="loading">Loading recipes...</div>;
  if (error) return <div className="error">Failed to load recipes: {error.message}</div>;

  return (
    <div className="management-page">
      <PageHeader
        title="Crafting Recipes"
        backTo="/dashboard"
        actions={
          <Link to="/recipes/new">
            <Button variant="primary">+ Add Recipe</Button>
          </Link>
        }
      />

      <DataTable
        columns={columns}
        data={recipes ?? []}
        getKey={(row) => row.name}
        emptyMessage="No recipes found. Create your first crafting recipe!"
      />

      {deletingRecipe && (
        <DeleteConfirmation
          recipe={deletingRecipe}
          onConfirm={handleDelete}
          onCancel={() => setDeletingRecipe(null)}
          isLoading={deleteMutation.isPending}
        />
      )}
    </div>
  );
}
