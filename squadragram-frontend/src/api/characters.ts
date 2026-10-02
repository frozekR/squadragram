import type {
	Character,
	CharacterFilter,
	CreateCharacterDTO,
} from "#/shared/types/character";
import { API_BASE_URL, fetchJSON } from "./client";

export async function fetchCharacters(
	filter?: CharacterFilter,
): Promise<Character[]> {
	const params = new URLSearchParams();
	if (filter?.role) {
		params.append("role", filter.role);
	}
	if (filter?.name) {
		params.append("name", filter.name);
	}

	const response = await fetch(
		`${API_BASE_URL}/characters?${params.toString()}`,
	);
	if (!response.ok) {
		throw new Error(`Failed to fetch characters: ${response.statusText}`);
	}
	return response.json();
}

export function fetchCharacterByUUID(uuid: string): Promise<Character> {
	return fetchJSON(`/characters/${encodeURIComponent(uuid)}`);
}

export async function createCharacter(
	character: CreateCharacterDTO,
): Promise<Character> {
	const response = await fetch(`${API_BASE_URL}/characters`, {
		method: "POST",
		headers: {
			"Content-Type": "application/json",
		},
		body: JSON.stringify(character),
	});
	if (!response.ok) {
		throw new Error(`Failed to create character: ${response.statusText}`);
	}
	return response.json();
}

export async function updateCharacter(
	id: number,
	character: Partial<Omit<Character, "id">>,
): Promise<Character> {
	const response = await fetch(`${API_BASE_URL}/characters/${id}`, {
		method: "PUT",
		headers: {
			"Content-Type": "application/json",
		},
		body: JSON.stringify(character),
	});
	if (!response.ok) {
		throw new Error(`Failed to update character: ${response.statusText}`);
	}
	return response.json();
}

export async function deleteCharacter(id: number): Promise<void> {
	const response = await fetch(`${API_BASE_URL}/characters/${id}`, {
		method: "DELETE",
	});
	if (!response.ok) {
		throw new Error(`Failed to delete character: ${response.statusText}`);
	}
}
