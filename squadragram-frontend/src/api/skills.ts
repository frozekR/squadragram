import { fetchJSON } from "./client";
import type { Skill } from "#/shared/types/skill";

export function fetchCharacterSkills(characterUUID: string): Promise<Skill[]> {
	return fetchJSON(`/characters/${encodeURIComponent(characterUUID)}/skills`);
}
