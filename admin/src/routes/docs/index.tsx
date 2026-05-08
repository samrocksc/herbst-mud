import { createFileRoute, Link } from '@tanstack/react-router'
import { PageHeader } from '../../components/PageHeader'
import { DocsIcon } from '../../components/icons/DocsIcon'

export const Route = createFileRoute('/docs/')({
  component: DocsIndex,
})

const DOC_PAGES = [
  {
    title: 'Ability System',
    path: '/docs/ability-system',
    desc: 'Combat abilities, effect types, scaling, costs, and the classless skill system.',
  },
  {
    title: 'Combat Guide',
    path: '/docs/combat-guide',
    desc: 'Tick-based combat flow, damage formula, dodge/parry, and skill combos.',
  },
  {
    title: 'Trainable Skills',
    path: '/docs/trainable-skills',
    desc: 'How weapon/magic proficiencies relate to abilities and training mechanics.',
  },
  {
    title: 'NPC System',
    path: '/docs/npc-system',
    desc: 'NPC lifecycle, templates, instances, respawn, level scaling, and race effects.',
  },
  {
    title: 'Item System',
    path: '/docs/item-system',
    desc: 'Equipment slots, damage/armor calculation, tags, and item categories.',
  },
  {
    title: 'Faction System',
    path: '/docs/faction-system',
    desc: 'Standing mechanics, categories, and how factions affect gameplay.',
  },
  {
    title: 'Examine Skill',
    path: '/docs/examine-skill',
    desc: 'The examine command, hidden details, skill levels, and DC checks.',
  },
  {
    title: 'Config Reference',
    path: '/docs/config-reference',
    desc: 'What each config key does and how it affects the game world.',
  },
]

function DocsIndex() {
  return (
    <div className="management-page">
      <PageHeader title="Documentation" backTo="/dashboard" />

      <p className="text-text-muted mb-6 max-w-2xl">
        Reference docs for Herbst MUD game mechanics. These pages explain how the
        systems behind the admin panel work in-game. Hover over form fields in the
        admin pages for inline tooltips, or read these pages for deeper context.
      </p>

      <div className="grid grid-cols-[repeat(auto-fit,minmax(280px,1fr))] gap-4">
        {DOC_PAGES.map((page) => (
          <Link
            key={page.path}
            to={page.path}
            className="block bg-surface-muted rounded-lg p-5 no-underline text-text border border-border transition-colors hover:border-primary hover:bg-surface-muted/70"
          >
            <div className="flex items-center gap-2 mb-2">
              <DocsIcon stroke="var(--color-primary)" />
              <span className="font-semibold">{page.title}</span>
            </div>
            <div className="text-text-muted text-sm">{page.desc}</div>
          </Link>
        ))}
      </div>
    </div>
  )
}