import { findMediaURL } from "#/api/media";
import { skillTypeLabels, type Skill } from "#/shared/types/skill";
import { DemoVideo } from "./DemoVideo";
import { MediaImage } from "./MediaImage";

export function SkillCard({ skill }: { skill: Skill }) {
	const icon = findMediaURL(skill.media_metadata, "ICON");
	const render = findMediaURL(skill.media_metadata, "RENDER");
	const demo = findMediaURL(skill.media_metadata, "DEMO");
	return (
		<article className="skill-card">
			<div className="skill-heading">
				<MediaImage
					src={icon}
					alt={`${skill.name} icon`}
					className="skill-icon"
					fallback="✦"
				/>
				<div>
					<p className="eyebrow">{skillTypeLabels[skill.type] || skill.type}</p>
					<h3>{skill.name}</h3>
				</div>
			</div>
			<p className="skill-description">{skill.description}</p>
			{render && (
				<MediaImage
					src={render}
					alt={`${skill.name} render`}
					className="skill-render"
					fit="contain"
				/>
			)}
			{demo ? (
				<div className="skill-demo">
					<h4>See it in action</h4>
					<DemoVideo src={demo} title={`${skill.name} demo`} poster={render} />
				</div>
			) : (
				<p className="media-note">Demo coming soon</p>
			)}
		</article>
	);
}
