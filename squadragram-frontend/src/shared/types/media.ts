export type MediaType = "ICON" | "RENDER" | "DEMO";
export type OwnerType = "CHARACTER" | "SKILL" | "EMOTE" | "SKIN";

export interface MediaMetadata {
	id: number;
	owner_type: OwnerType;
	owner_uuid: string;
	media_type: MediaType;
	object_key: string;
}
