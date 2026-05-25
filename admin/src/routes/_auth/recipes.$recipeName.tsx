import { createFileRoute, useNavigate, useLocation, Outlet } from "@tanstack/react-router";
import { Button } from "../../components/Button";
import { PageHeader } from "../../components/PageHeader";
import { PageContainer } from "../../components/PageContainer";
import { useRecipe } from "../../hooks/useRecipes";

export const Route = createFileRoute("/_auth/recipes/$recipeName")({
  component: RecipeDetail,
});

function RecipeDetail() {
  const { recipeName } = Route.useParams();
  const navigate = useNavigate();
  const location = useLocation();

  const recipeQuery = useRecipe(recipeName);
  const recipe = recipeQuery.data;

  // If we're on the edit child route, render outlet so the child route handles it
  if (location.pathname.endsWith("/edit")) {
    return <Outlet />;
  }

  if (recipeQuery.isLoading) return <div className="loading">Loading recipe...</div>;
  if (recipeQuery.error || !recipe) return <div className="error">Failed to load recipe.</div>;

  return (
    <PageContainer>
      <PageHeader
        title={recipe.display_name || recipe.name}
        backTo="/recipes"
        actions={
          <Button
            variant="primary"
            size="sm"
            onClick={() => navigate({ to: "/recipes/$recipeName/edit", params: { recipeName } })}
          >
            Edit
          </Button>
        }
      />

      <div className="card bg-surface p-6 border border-border rounded space-y-3">
        <div className="grid grid-cols-2 gap-4 text-sm">
          <div><span className="text-text-muted">Name:</span> <span className="text-text font-medium">{recipe.name}</span></div>
          <div><span className="text-text-muted">Display Name:</span> <span className="text-text font-medium">{recipe.display_name}</span></div>
          <div><span className="text-text-muted">Station:</span> <span className="text-text font-medium">{recipe.required_station_tag}</span></div>
          <div><span className="text-text-muted">Class:</span> <span className="text-text font-medium">{recipe.required_class || "Any"}</span></div>
          <div><span className="text-text-muted">Skill:</span> <span className="text-text font-medium">{recipe.required_skill || "—"}</span></div>
          <div><span className="text-text-muted">Skill Level:</span> <span className="text-text font-medium">{recipe.required_skill_level}</span></div>
          <div><span className="text-text-muted">Craft Time:</span> <span className="text-text font-medium">{recipe.craft_time_secs}s</span></div>
          <div><span className="text-text-muted">World:</span> <span className="text-text font-medium">{recipe.world_id}</span></div>
        </div>

        {recipe.description && (
          <div className="text-sm text-text-muted">{recipe.description}</div>
        )}

        <div className="space-y-2">
          <h4 className="text-text font-semibold text-sm mt-4">Inputs</h4>
          {(recipe.inputs?.length ?? 0) === 0 ? (
            <div className="text-text-muted text-sm">No inputs.</div>
          ) : (
            <ul className="list-disc pl-5 text-sm text-text space-y-1">
              {recipe.inputs.map((input, i) => (
                <li key={i}>
                  {input.quantity}x {input.equipment_template_slug || "(unset)"}
                  {input.consumed ? " (consumed)" : ""}
                </li>
              ))}
            </ul>
          )}
        </div>

        <div className="space-y-2">
          <h4 className="text-text font-semibold text-sm mt-4">Outputs</h4>
          {(recipe.outputs?.length ?? 0) === 0 ? (
            <div className="text-text-muted text-sm">No outputs.</div>
          ) : (
            <ul className="list-disc pl-5 text-sm text-text space-y-1">
              {recipe.outputs.map((output, i) => (
                <li key={i}>
                  {output.quantity}x {output.equipment_template_slug || "(unset)"}
                </li>
              ))}
            </ul>
          )}
        </div>
      </div>
    </PageContainer>
  );
}
