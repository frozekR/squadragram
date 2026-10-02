import { useState } from "react";
import { createFileRoute, Link } from "@tanstack/react-router";
import { findMediaURL } from "#/api/media";
import { MediaImage } from "#/components/MediaImage";
import { SkillSelector } from "#/components/SkillSelector";
import { useCharacter } from "#/hooks/useCharacter";
import { useCharacterSkills } from "#/hooks/useCharacterSkills";

export const Route = createFileRoute("/characters/$characterId")({
	component: HeroPage,
});

function HeroPage() {
	const [activeTab, setActiveTab] = useState("skills");
	// Resolve the sequential ID from the URL before requesting skills by UUID.
	const { characterId } = Route.useParams();
	const characterQuery = useCharacter(characterId);
	const skillsQuery = useCharacterSkills(characterQuery.data?.uuid);
	const character = characterQuery.data;

	if (characterQuery.isLoading)
		return (
			<main className="hero-page content-status" aria-live="polite">
				Loading hero…
			</main>
		);
	if (characterQuery.isError || !character)
		return (
			<main className="hero-page content-status" role="alert">
				<h1>Could not load this hero.</h1>
				<p>Check the link and try again.</p>
				<button type="button" onClick={() => void characterQuery.refetch()}>
					Try again
				</button>
				<Link to="/characters">Back to heroes</Link>
			</main>
		);

	const icon = findMediaURL(character.media_metadata, "ICON");
	const render = findMediaURL(character.media_metadata, "RENDER");
	const skills = skillsQuery.data || [];
	const tabs = [
		{
			id: "skills",
			label: `Скиллы`,
		},
		{ id: "skins", label: "Скины" },
		{ id: "emotes", label: "Эмоции" },
	];

	return (
		<main className="hero-page">
			<Link to="/characters" className="back-link">
				← All heroes
			</Link>
			<section className={`hero-banner role-${character.role.toLowerCase()}`}>
				<div className="hero-intro">
					<p className="eyebrow">Профиль героя</p>
					<h1>{character.name}</h1>
					<span className="role-badge">{character.role}</span>
					<p>{character.description}</p>
				</div>
				<MediaImage
					key={character.uuid}
					src={render || icon}
					alt={character.name}
					className="hero-art"
					fit="contain"
					loading="eager"
					fallback="Арт Героя недоступен"
				/>
			</section>
			<nav className="hero-tabs" aria-label="Hero sections">
				{tabs.map((tab) => (
					<button
						key={tab.id}
						type="button"
						aria-pressed={activeTab === tab.id}
						onClick={() => setActiveTab(tab.id)}
					>
						{tab.label}
					</button>
				))}
			</nav>
			<section
				className="hero-content"
				aria-label={tabs.find((tab) => tab.id === activeTab)?.label}
			>
				{activeTab === "skills" && (
					<>
						<div className="section-heading">
							<div>
								<p className="eyebrow">Moveset</p>
								<h2>Скиллы & Атаки</h2>
							</div>
							<button
								type="button"
								className="refresh-button"
								disabled={skillsQuery.isFetching}
								onClick={() => void skillsQuery.refetch()}
							>
								{skillsQuery.isFetching ? "Обновляется..." : "Обновить"}
							</button>
						</div>
						{skillsQuery.isLoading ? (
							<p className="content-status" aria-live="polite">
								Loading skills…
							</p>
						) : skillsQuery.isError ? (
							<div className="content-status" role="alert">
								<p>Could not load skills.</p>
								<button
									type="button"
									onClick={() => void skillsQuery.refetch()}
								>
									Try again
								</button>
							</div>
						) : skills.length === 0 ? (
							<p className="content-status">This hero has no skills yet.</p>
						) : (
							<SkillSelector key={character.uuid} skills={skills} />
						)}
					</>
				)}
				{activeTab === "skins" && (
					<div className="content-status">
						<h2>Skins</h2>
						<p>No skins available yet.</p>
					</div>
				)}
				{activeTab === "emotes" && (
					<div className="content-status">
						<h2>Emotes</h2>
						<p>No emotes available yet.</p>
					</div>
				)}
			</section>
		</main>
	);
}
