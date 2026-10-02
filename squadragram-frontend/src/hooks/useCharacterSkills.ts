import { useQuery } from "@tanstack/react-query";
import { fetchCharacterSkills } from "#/api/skills";

export function useCharacterSkills(uuid = "") {
	return useQuery({
		queryKey: ["character", uuid, "skills"],
		queryFn: () => fetchCharacterSkills(uuid),
		enabled: uuid.length > 0,
	});
}
