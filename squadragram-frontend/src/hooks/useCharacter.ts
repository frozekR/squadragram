import { useQuery } from '@tanstack/react-query'

import { fetchCharacterById } from '#/api/characters'

export function useCharacter(id: number) {
  return useQuery({
    queryKey: ['character', id],
    queryFn: () => fetchCharacterById(id),
  })
}