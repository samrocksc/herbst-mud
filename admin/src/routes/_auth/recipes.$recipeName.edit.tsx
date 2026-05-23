import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { Button } from "../../components/Button";
import { PageHeader } from "../../components/PageHeader";
import { useRecipe, useUpdateRecipe, type RecipeInput } from "../../hooks/useRecipes";
import { useWorldStore } from "../../contexts/WorldStoreContext";
import { RecipeForm } from "./RecipeForm";

export const Route = createFileRoute("/_auth/recipes/$recipeName/edit")({
  component: RecipeEditPage,
});

function RecipeEditPage() {
  const { recipeName } = Route.useParams();
  const navigate = useNavigate();
  const { currentWorld } = useWorldStore();

  const recipeQuery = useRecipe(recipeName);
  const recipe = recipeQuery.data;
  const updateMutation = useUpdateRecipe();

  const handleSubmit = async (input: RecipeInput) => {
    try {
      await updateMutation.mutateAsync({ name: recipeName, input });
      navigate({ to: "/recipes/$recipeName", params: { recipeName } });
    } catch { /* error handled by mutation state */ }
  };

  if (recipeQuery.isLoading) return <div className="loading">Loading recipe...</div>;
  if (recipeQuery.error || !recipe) return <div className="error">Failed to load recipe.</div>;

  return (
    <div className="management-page">
      <PageHeader
        title={`Edit: ${recipe.display_name || recipe.name}`}
        backTo="/recipes"
        actions={
          <Button variant="secondary" size="sm" onClick={() => navigate({ to: "/recipes/$recipeName", params: { recipeName } })} >
            Cancel
          </Button>
        }
      />

      <RecipeForm
        recipe={recipe}
        onSubmit={handleSubmit}
        onCancel={() => navigate({ to: "/recipes/$recipeName", params: { recipeName } })}
      />
    </div>
  );
}
