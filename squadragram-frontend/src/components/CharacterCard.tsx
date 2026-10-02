import { Link } from "@tanstack/react-router";
import { findMediaURL } from "#/api/media";
import type { Character } from "#/shared/types/character";
import { MediaImage } from "./MediaImage";

export default function CharacterCard({ character }: { character: Character }) {
	const image =
		findMediaURL(character.media_metadata, "ICON") ||
		findMediaURL(character.media_metadata, "RENDER");
	return (
		<Link
			to="/characters/$characterId"
			params={{ characterId: String(character.id) }}
			className={`character-card role-${character.role.toLowerCase()}`}
		>
			<MediaImage
				src={image}
				alt={character.name}
				className="character-card-image"
				fallback={character.name
					.split(" ")
					.map((word) => word[0])
					.slice(0, 3)
					.join("")}
			/>
			<span className="role-badge">{character.role}</span>
			<div className="character-card-caption">
				<h3>{character.name}</h3>
				<span>
					View hero <span aria-hidden="true">↗</span>
				</span>
			</div>
		</Link>
	);
}
