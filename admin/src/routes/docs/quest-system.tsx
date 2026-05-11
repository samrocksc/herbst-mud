import { createFileRoute } from '@tanstack/react-router'
import { PageHeader } from '../../components/PageHeader'

export const Route = createFileRoute('/docs/quest-system')({
  component: QuestSystemDoc,
})

function Section({ title, children }: Readonly<{ title: string; children: React.ReactNode }>) {
  return (
    <section className="mb-8">
      <h2 className="text-lg font-semibold text-text mb-3 pb-2 border-b border-border">{title}</h2>
      {children}
    </section>
  )
}

function InfoBox({ children }: Readonly<{ children: React.ReactNode }>) {
  return (
    <div className="bg-primary/10 border border-primary/30 rounded-lg p-4 mb-4 text-sm">
      {children}
    </div>
  )
}

function Table({
  headers,
  rows,
}: Readonly<{ headers: string[]; rows: (string | React.ReactNode)[][] }>) {
  return (
    <div className="overflow-x-auto mb-4">
      <table className="w-full text-sm border border-border rounded-lg">
        <thead>
          <tr className="bg-surface-muted">
            {headers.map((h) => (
              <th key={h} className="text-left px-3 py-2 font-semibold border-b border-border">{h}</th>
            ))}
          </tr>
        </thead>
        <tbody>
          {rows.map((row, i) => (
            <tr key={i} className="border-b border-border last:border-0">
              {row.map((cell, j) => (
                <td key={j} className="px-3 py-2">{cell}</td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

function QuestSystemDoc() {
  return (
    <div className="management-page max-w-4xl">
      <PageHeader title="Quest System" backTo="/docs" />

      <InfoBox>
        <strong>TL;DR:</strong> Quests are objective-based tasks that characters accept, progress through
        sequentially, and complete for rewards. Objectives track kills, exploration, item collection,
        NPC conversations, and item returns. Progress is tracked per-character with automatic advancement
        from in-game events.
      </InfoBox>

      <Section title="Quest Lifecycle">
        <p className="text-text-muted mb-3">
          A quest progresses through these states for each character:
        </p>
        <Table
          headers={['State', 'Description']}
          rows={[
            ['Available', 'Quest meets prerequisites and character hasn\'t accepted it yet.'],
            ['Active', 'Character has accepted the quest and is working on objectives.'],
            ['Completed', 'All objectives finished. Rewards applied.'],
            ['Abandoned', 'Character chose to abandon. Progress is lost.'],
          ]}
        />
      </Section>

      <Section title="Objective Types">
        <p className="text-text-muted mb-3">
          Each quest has an ordered list of objectives. Objectives are completed
          sequentially — later objectives don't start counting until earlier ones are done.
        </p>
        <Table
          headers={['Type', 'Target', 'Description']}
          rows={[
            ['kill', 'NPC template ID', 'Defeat a specific NPC type. Count tracks kills.'],
            ['explore', 'Room ID', 'Visit a specific room. Auto-completes on room entry.'],
            ['collect', 'Item template ID', 'Gather items. Count tracks items picked up.'],
            ['talk', 'NPC template ID', 'Speak with a specific NPC. Auto-completes on conversation.'],
            ['return', 'Item template ID', 'Bring an item back to a quest giver.'],
          ]}
        />
      </Section>

      <Section title="Repeat Modes">
        <p className="text-text-muted mb-3">
          Quests can be one-time or repeatable:
        </p>
        <Table
          headers={['Mode', 'Behavior']}
          rows={[
            ['none', 'One-time quest. Cannot be re-accepted after completion.'],
            ['cooldown', 'Repeatable after cooldown_hours since last completion.'],
            ['always', 'Always repeatable, no cooldown enforced.'],
          ]}
        />
        <p className="text-text-muted">
          Prerequisites are checked before acceptance. A character cannot accept a quest
          they already have active.
        </p>
      </Section>

      <Section title="Rewards">
        <p className="text-text-muted mb-3">
          Completing a quest grants these reward types:
        </p>
        <Table
          headers={['Reward', 'Description']}
          rows={[
            ['XP', 'Experience points added to the character.'],
            ['Items', 'Item templates granted to character inventory.'],
            ['Effects', 'Ability effects applied to the character.'],
            ['Tags', 'Tags added or removed from the character.'],
            ['Achievements', 'Achievements unlocked for the character.'],
          ]}
        />
        <InfoBox>
          The admin UI currently only exposes the XP reward field. Item, effect, tag,
          and achievement rewards can be set via the API directly.
        </InfoBox>
      </Section>

      <Section title="Player Commands">
        <Table
          headers={['Command', 'Description']}
          rows={[
            ['quests / quest / q', 'Show quest tracker with all active/completed/abandoned quests.'],
            ['quest accept &lt;id&gt;', 'Accept a quest by its ID. Checks prerequisites.'],
            ['quest abandon &lt;id&gt;', 'Abandon an active quest. Progress is lost.'],
          ]}
        />
        <p className="text-text-muted">
          Quest progress advances automatically when in-game events match objective
          types (killing NPCs, entering rooms, picking up items, talking to NPCs).
        </p>
      </Section>

      <Section title="Admin API">
        <Table
          headers={['Method', 'Endpoint', 'Description']}
          rows={[
            ['GET', '/api/quests', 'List all quest definitions.'],
            ['POST', '/api/quests', 'Create a new quest.'],
            ['GET', '/api/quests/:id', 'Get a quest definition by ID.'],
            ['PUT', '/api/quests/:id', 'Update a quest definition.'],
            ['DELETE', '/api/quests/:id', 'Delete a quest (fails if progress records exist).'],
            ['GET', '/api/characters/:id/quests', 'List quest progress for a character.'],
            ['POST', '/api/characters/:id/quests', 'Accept a quest (body: {"quest_id": N}).'],
            ['PUT', '/api/characters/:id/quests/:qid/check', 'Check/increment single quest progress.'],
            ['PUT', '/api/characters/:id/quests/:qid/abandon', 'Abandon a quest.'],
            ['POST', '/api/characters/:id/quests/check-all', 'Bulk check all matching quests.'],
          ]}
        />
      </Section>
    </div>
  )
}