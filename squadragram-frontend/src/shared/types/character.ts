import type { MediaMetadata } from "./media";

export type CharacterRole =
	| "DAMAGE"
	| "TANK"
	| "TECHNICAL"
	| "MELEE"
	| "RANGED";

export interface Character {
	id: number;
	uuid: string;
	name: string;
	description: string;
	role: CharacterRole;
	media_metadata: MediaMetadata[];
}

export type CreateCharacterDTO = Pick<
	Character,
	"name" | "description" | "role"
>;

export interface CharacterFilter {
	role?: CharacterRole;
	name?: string;
}
