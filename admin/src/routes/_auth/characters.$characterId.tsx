import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { apiGet } from '../../utils/apiFetch'
import { PageHeader } from '../../components/PageHeader'
import { Button } from '../../components/Button'
import { EquippedItemsView } from '../../components/EquippedItemsView'
import { CharacterStats } from './-characters.$characterId.stats'

export const Route = createFileRoute('/_auth/characters/$characterId')({
  component: CharacterDetail,
})

type Character = Readonly<{
  id: number
  name: string
  isNPC: boolean
  currentRoomId: number
  hitpoints: number
  max_hitpoints: number
  stamina: number
  max_stamina: number
  mana: number
  max_mana: number
  race: string
  class: string
  level: number
  xp: number
  strength: number
  dexterity: number
  intelligence: number
  wisdom: number
  constitution: number
  gender: string
  description: string
  is_admin: boolean
  is_immortal: boolean
}>

function CharacterDetail() {
  const { characterId } = Route.useParams()
  const [showEquipped, setShowEquipped] = useState(true)

  const { data: character, isLoading, error } = useQuery<Character>({
    queryKey: ['character', characterId],
    queryFn: () => apiGet<Character>(`${window.location.origin}/characters/${characterId}`),
  })

  if (isLoading) return <div className="p-8"><PageHeader title="Loading..." backTo="/players" /></div>
  if (error || !character) return <div className="p-8"><PageHeader title="Error" backTo="/players" /><div className="text-danger">Failed to load character</div></div>

  return (
    <div className="p-8">
      <PageHeader title={character.name} backTo="/players" actions={
        <Button variant={showEquipped ? 'primary' : 'secondary'} size="sm" onClick={() => setShowEquipped(!showEquipped)}>
          {showEquipped ? 'Hide Equipment' : 'Show Equipment'}
        </Button>
      } />
      <div className="max-w-3xl">
        <CharacterStats character={character} />
        {showEquipped && <EquippedItemsView characterId={character.id} characterRace={character.race} />}
      </div>
    </div>
  )
}