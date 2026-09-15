import { useQuery } from '@tanstack/react-query'

import { fetchCharacters } from '#/api/characters'

export function useCharacters() {
  return useQuery({
    queryKey: ['characters'],
    queryFn: () => fetchCharacters(),
  })
}