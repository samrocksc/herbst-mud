import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/_auth/skills')({
  component: SkillsPage,
})

function SkillsPage() {
  return (
    <div className="page-header">
      <h1>⭐ Skills & Talents</h1>
      <p>Manage character skills and talents for all classes.</p>
      <div className="placeholder-page">
        <h2>Coming Soon</h2>
        <p>Skill and talent management features will be available here.</p>
        <p>Features include: Add skills/talents, Edit properties, Assign to classes, Set requirements</p>
      </div>
    </div>
  )
}