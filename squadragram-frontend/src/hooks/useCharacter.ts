import { useQuery, useQueryClient } from "@tanstack/react-query";

import { fetchCharacterByUUID, fetchCharacters } from "#/api/characters";

export function useCharacter(id: string) {
	const queryClient = useQueryClient();
	return useQuery({
		queryKey: ["character-by-id", id],
		queryFn: async () => {
			if (!/^[1-9]\d*$/.test(id) || !Number.isSafeInteger(Number(id))) {
				throw new Error("Invalid character ID");
			}
			// The browser URL uses the sequential ID. The API still uses UUIDs.
			// Reuse the heroes list cache, including when opening a direct link.
			const characters = await queryClient.fetchQuery({
				queryKey: ["characters"],
				queryFn: () => fetchCharacters(),
			});
			const character = characters.find((item) => item.id === Number(id));
			if (!character) throw new Error("Character not found");
			return fetchCharacterByUUID(character.uuid);
		},
	});
}
