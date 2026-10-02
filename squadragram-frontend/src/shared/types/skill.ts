import type { MediaMetadata } from "./media";

export type SkillType =
	| "PASSIVE"
	| "RUSH_ATTACK"
	| "SKILL"
	| "SUPER_ATTACK"
	| "MAX_SUPER_ATTACK"
	| "TRANSFORMATION";

export interface Skill {
	id: number;
	uuid: string;
	name: string;
	description: string;
	type: SkillType;
	sort_order: number;
	character_id: number | null;
	media_metadata: MediaMetadata[];
}

export const skillTypeLabels: Record<SkillType, string> = {
	PASSIVE: "Passive",
	RUSH_ATTACK: "Rush attack",
	SKILL: "Skill",
	SUPER_ATTACK: "Super attack",
	MAX_SUPER_ATTACK: "Max super attack",
	TRANSFORMATION: "Transformation",
};

export const skillCategoryOrder: SkillType[] = [
	"PASSIVE",
	"RUSH_ATTACK",
	"SKILL",
	"SUPER_ATTACK",
	"MAX_SUPER_ATTACK",
	"TRANSFORMATION",
];

export function groupSkills(skills: Skill[]) {
	return skillCategoryOrder
		.map((type) => ({
			type,
			skills: skills
				.filter((skill) => skill.type === type)
				.sort(
					(a, b) => (a.sort_order ?? 0) - (b.sort_order ?? 0) || a.id - b.id,
				),
		}))
		.filter((group) => group.skills.length > 0);
}
