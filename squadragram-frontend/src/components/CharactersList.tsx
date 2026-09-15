import CharacterCard from './CharacterCard'
import { useCharacters } from '#/hooks/useCharacters'

type Character = {
  id: string | number
  name: string
  role: string
}

export const CharactersList = () => {
  const { data: characters = [] as Character[], isLoading, isError } = useCharacters() as {
    data?: Character[]
    isLoading: boolean
    isError: boolean
  }

  if (isLoading) {
    return <p>Loading characters...</p>
  }

  if (isError) {
    return <p>Failed to load characters.</p>
  }

  return (
    <div className="flex flex-wrap gap-4">
      {characters.map((character) => (
        <CharacterCard
          key={character.id}
          id={character.id}
          name={character.name}
          role={character.role}
          imageUrl={`/placeholder_media/roles/icon_${character.role.toLowerCase()}.png`}
        />
      ))}
    </div>
  )
}