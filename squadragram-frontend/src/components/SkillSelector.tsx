import { useId, useState } from "react";
import { findMediaURL } from "#/api/media";
import { groupSkills, skillTypeLabels, type Skill } from "#/shared/types/skill";
import { MediaImage } from "./MediaImage";
import { SkillCard } from "./SkillCard";

export function SkillSelector({ skills }: { skills: Skill[] }) {
	const [selectedUUID, setSelectedUUID] = useState("");
	const detailID = useId();
	const groups = groupSkills(skills);
	// Preserve the choice on refresh, and fall back if that skill was moved or removed.
	const selected =
		skills.find((skill) => skill.uuid === selectedUUID) || groups[0]?.skills[0];
	if (!selected) return null;
	return (
		<div className="skill-selector">
			<nav className="skill-category-strip" aria-label="Choose a skill">
				{groups.map((group) => (
					<div className="skill-category" key={group.type}>
						<p className="eyebrow">{skillTypeLabels[group.type]}</p>
						<div className="skill-category-icons">
							{group.skills.map((skill) => (
								<button
									key={skill.uuid}
									type="button"
									className="skill-choice"
									aria-label={skill.name}
									aria-pressed={selected.uuid === skill.uuid}
									aria-controls={detailID}
									title={skill.name}
									onClick={() => setSelectedUUID(skill.uuid)}
								>
									<MediaImage
										src={findMediaURL(skill.media_metadata, "ICON")}
										alt=""
										className="skill-choice-icon"
										fallback={skill.name.slice(0, 1)}
									/>
								</button>
							))}
						</div>
					</div>
				))}
			</nav>
			<div id={detailID} className="skill-selected-detail">
				{/* Remount to stop the previous demo and reset its playback/error state. */}
				<SkillCard key={selected.uuid} skill={selected} />
			</div>
		</div>
	);
}
