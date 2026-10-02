import { adminUI as ui } from "#/components/admin/styles";
import { useCallback, useEffect, useState } from "react";
import { createFileRoute, Link } from "@tanstack/react-router";
import { useRef } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { AdminLogin } from "#/components/admin/AdminLogin";
import {
	HeroEditor,
	SkillEditor,
	roleLabels,
} from "#/components/admin/EntityEditor";
import { MediaManager } from "#/components/admin/MediaManager";
import { MediaImage } from "#/components/MediaImage";
import { findMediaURL } from "#/api/media";
import { useCharacters } from "#/hooks/useCharacters";
import { useCharacterSkills } from "#/hooks/useCharacterSkills";
import type { Character } from "#/shared/types/character";
import { groupSkills, skillTypeLabels, type Skill } from "#/shared/types/skill";

export const Route = createFileRoute("/admin/")({ component: AdminPage });

function AdminPage() {
	// The secret is entered by the administrator and lives only in this component.
	// Never persist it in localStorage, query caches, URLs or the frontend bundle.
	const [token, setToken] = useState("");
	return token ? (
		<AdminWorkspace token={token} onLogout={() => setToken("")} />
	) : (
		<AdminLogin onLogin={setToken} />
	);
}

function AdminWorkspace({
	token,
	onLogout,
}: {
	token: string;
	onLogout: () => void;
}) {
	const queryClient = useQueryClient();
	const heroesQuery = useCharacters();
	const heroes = heroesQuery.data || [];
	const [selectedUUID, setSelectedUUID] = useState("");
	const [newHero, setNewHero] = useState(false);
	const [tab, setTab] = useState<"hero" | "skills">("hero");
	const [search, setSearch] = useState("");
	const [selectedSkillUUID, setSelectedSkillUUID] = useState("");
	const [newSkill, setNewSkill] = useState(false);
	const [editorDirty, setDirty] = useState(false);
	const [mediaDirty, setMediaDirty] = useState<Record<string, boolean>>({});
	const dirty = editorDirty || Object.values(mediaDirty).some(Boolean);
	const onMediaDirty = useCallback((key: string, value: boolean) => {
		setMediaDirty((old) =>
			old[key] === value ? old : { ...old, [key]: value },
		);
	}, []);
	const [notice, setNotice] = useState("");
	const [pendingNavigation, setPendingNavigation] = useState<
		(() => void) | null
	>(null);
	const onDirty = useCallback((value: boolean) => setDirty(value), []);
	const hero = heroes.find((item) => item.uuid === selectedUUID);
	const skillsQuery = useCharacterSkills(hero?.uuid);
	const skills = skillsQuery.data || [];
	const skill = skills.find((item) => item.uuid === selectedSkillUUID);
	useEffect(() => {
		if (!dirty) return;
		const warn = (event: BeforeUnloadEvent) => {
			event.preventDefault();
			event.returnValue = "";
		};
		window.addEventListener("beforeunload", warn);
		return () => window.removeEventListener("beforeunload", warn);
	}, [dirty]);
	function navigate(action: () => void) {
		if (dirty) {
			setPendingNavigation(() => action);
			return;
		}
		setDirty(false);
		setNotice("");
		action();
	}
	async function refresh() {
		await Promise.all([
			queryClient.invalidateQueries({ queryKey: ["characters"] }),
			queryClient.invalidateQueries({ queryKey: ["character"] }),
			queryClient.invalidateQueries({ queryKey: ["character-by-id"] }),
		]);
	}
	function heroSaved(saved: Character) {
		queryClient.setQueryData<Character[]>(["characters"], (old) =>
			old?.some((item) => item.uuid === saved.uuid)
				? old.map((item) => (item.uuid === saved.uuid ? saved : item))
				: [...(old || []), saved],
		);
		setSelectedUUID(saved.uuid);
		setNewHero(false);
		setDirty(false);
		setNotice(
			newHero
				? "Герой создан. Теперь можно добавить медиа и скиллы."
				: "Изменения героя сохранены.",
		);
		void refresh();
	}
	function skillSaved(saved: Skill, ownerUUID: string) {
		queryClient.setQueryData<Skill[]>(
			["character", ownerUUID, "skills"],
			(old) =>
				old?.some((item) => item.uuid === saved.uuid)
					? old.map((item) => (item.uuid === saved.uuid ? saved : item))
					: [...(old || []), saved],
		);
		setSelectedUUID(ownerUUID);
		setSelectedSkillUUID(saved.uuid);
		setNewSkill(false);
		setDirty(false);
		setNotice(
			newSkill
				? "Скилл создан. Можно загрузить его медиа."
				: "Изменения скилла сохранены.",
		);
		void refresh();
	}
	const visibleHeroes = heroes.filter((item) =>
		item.name.toLowerCase().includes(search.toLowerCase()),
	);
	return (
		<main className={ui.workspace}>
			<div className={ui.topbar}>
				<div>
					<p className={ui.eyebrow}>Squadragram · управление</p>
					<h1>Админка</h1>
				</div>
				<div className={ui.actions}>
					<span className={ui.access}>
						<span aria-hidden="true">●</span> Доступ открыт
					</span>
					<button
						type="button"
						className={ui.secondary}
						onClick={() => navigate(onLogout)}
					>
						Выйти
					</button>
				</div>
			</div>
			<div className={ui.layout}>
				<aside className={ui.sidebar}>
					<div className={ui.sidebarHeading}>
						<h2>
							Герои <span>{heroes.length}</span>
						</h2>
						<button
							type="button"
							className={ui.primary}
							onClick={() =>
								navigate(() => {
									setNewHero(true);
									setSelectedUUID("");
									setTab("hero");
								})
							}
						>
							+ Добавить
						</button>
					</div>
					<label className={ui.search}>
						Поиск героев
						<input
							type="search"
							placeholder="Найти по имени…"
							value={search}
							onChange={(event) => setSearch(event.target.value)}
						/>
					</label>
					{heroesQuery.isLoading ? (
						<p className={ui.status}>Загружаем героев…</p>
					) : heroesQuery.isError ? (
						<div className={ui.error} role="alert">
							Не удалось загрузить героев.
							<button type="button" onClick={() => void heroesQuery.refetch()}>
								Повторить
							</button>
						</div>
					) : (
						<div className={ui.heroList}>
							{visibleHeroes.map((item) => (
								<button
									key={item.uuid}
									type="button"
									className={ui.heroItem}
									aria-pressed={selectedUUID === item.uuid}
									onClick={() => {
										if (selectedUUID === item.uuid) return;
										navigate(() => {
											setSelectedUUID(item.uuid);
											setNewHero(false);
											setSelectedSkillUUID("");
											setNewSkill(false);
										});
									}}
								>
									<MediaImage
										src={findMediaURL(item.media_metadata, "ICON")}
										alt={item.name}
										className={ui.avatar}
										fallback={item.name.slice(0, 1)}
									/>
									<span>
										<strong>{item.name}</strong>
										<small>{roleLabels[item.role]}</small>
									</span>
								</button>
							))}
							{visibleHeroes.length === 0 && (
								<p className={ui.status}>
									{search ? "Ничего не найдено." : "Добавьте первого героя."}
								</p>
							)}
						</div>
					)}
				</aside>
				<div className={ui.main}>
					{notice && (
						<div className={ui.banner} aria-live="polite">
							{notice}
							<button
								type="button"
								aria-label="Закрыть уведомление"
								onClick={() => setNotice("")}
							>
								×
							</button>
						</div>
					)}
					{newHero ? (
						<>
							<HeroEditor
								key="new-hero"
								token={token}
								onSaved={heroSaved}
								onDirty={onDirty}
							/>
							<div className={`${ui.panel} ${ui.help}`}>
								После создания героя здесь появится управление его медиа и
								скиллами.
							</div>
						</>
					) : hero ? (
						<>
							<div className={ui.entityHeading}>
								<div>
									<h2>{hero.name}</h2>
									<p>
										{roleLabels[hero.role]} · Герой #{hero.id}
									</p>
								</div>
								<Link
									to="/characters/$characterId"
									params={{ characterId: String(hero.id) }}
									target="_blank"
									rel="noopener noreferrer"
									className={ui.secondary}
								>
									Открыть страницу ↗
								</Link>
							</div>
							<nav className={ui.tabs} aria-label="Разделы редактирования">
								<button
									type="button"
									aria-pressed={tab === "hero"}
									onClick={() => {
										if (tab !== "hero") navigate(() => setTab("hero"));
									}}
								>
									Герой и медиа
								</button>
								<button
									type="button"
									aria-pressed={tab === "skills"}
									onClick={() => {
										if (tab !== "skills") navigate(() => setTab("skills"));
									}}
								>
									Скиллы {skillsQuery.data ? `· ${skills.length}` : ""}
								</button>
							</nav>
							{tab === "hero" ? (
								<>
									<HeroEditor
										key={hero.uuid}
										token={token}
										character={hero}
										onSaved={heroSaved}
										onDirty={onDirty}
									/>
									<MediaManager
										token={token}
										ownerType="CHARACTER"
										ownerUUID={hero.uuid}
										metadata={hero.media_metadata || []}
										onChanged={refresh}
										onDirty={onMediaDirty}
									/>
								</>
							) : (
								<>
									<section className={ui.panel}>
										<div className={ui.sectionHeading}>
											<div>
												<p className={ui.eyebrow}>Умения и атаки</p>
												<h2>Скиллы героя</h2>
											</div>
											<button
												type="button"
												className={ui.primary}
												onClick={() =>
													navigate(() => {
														setNewSkill(true);
														setSelectedSkillUUID("");
													})
												}
											>
												+ Добавить скилл
											</button>
										</div>
										{skillsQuery.isLoading ? (
											<p className={ui.status}>Загружаем скиллы…</p>
										) : skillsQuery.isError ? (
											<div className={ui.error} role="alert">
												Не удалось загрузить скиллы.
												<button
													type="button"
													onClick={() => void skillsQuery.refetch()}
												>
													Повторить
												</button>
											</div>
										) : (
											<div className={ui.skillList}>
												{groupSkills(skills).map((group) => (
													<div key={group.type} className={ui.skillCategory}>
														<p className={ui.eyebrow}>
															{skillTypeLabels[group.type]}
														</p>
														{group.skills.map((item) => (
															<button
																key={item.uuid}
																type="button"
																aria-pressed={selectedSkillUUID === item.uuid}
																onClick={() => {
																	if (
																		!newSkill &&
																		selectedSkillUUID === item.uuid
																	)
																		return;
																	navigate(() => {
																		setSelectedSkillUUID(item.uuid);
																		setNewSkill(false);
																	});
																}}
															>
																<MediaImage
																	src={findMediaURL(
																		item.media_metadata,
																		"ICON",
																	)}
																	alt={item.name}
																	className={ui.avatar}
																	fallback="✦"
																/>
																<span>
																	<strong>{item.name}</strong>
																	<small>Порядок: {item.sort_order ?? 0}</small>
																</span>
																<span aria-hidden="true">→</span>
															</button>
														))}
													</div>
												))}
												{skills.length === 0 && (
													<p className={ui.status}>У героя пока нет скиллов.</p>
												)}
											</div>
										)}
									</section>
									{(newSkill || skill) && (
										<>
											<SkillEditor
												key={newSkill ? `new-${hero.uuid}` : skill?.uuid}
												token={token}
												skill={newSkill ? undefined : skill}
												owner={hero}
												characters={heroes}
												defaultOrder={Math.min(
													2147483647,
													Math.max(
														-1,
														...skills.map((item) => item.sort_order ?? 0),
													) + 1,
												)}
												onSaved={skillSaved}
												onDirty={onDirty}
											/>
											{!newSkill && skill && (
												<MediaManager
													token={token}
													ownerType="SKILL"
													ownerUUID={skill.uuid}
													metadata={skill.media_metadata || []}
													onChanged={refresh}
													onDirty={onMediaDirty}
												/>
											)}
										</>
									)}
								</>
							)}
						</>
					) : (
						<div className={`${ui.welcome} ${ui.panel}`}>
							<span className={ui.mark} aria-hidden="true">
								S
							</span>
							<p className={ui.eyebrow}>Рабочее пространство</p>
							<h2>Всё начинается с героя</h2>
							<p>
								Выберите героя слева, чтобы изменить его данные, скиллы и медиа.
								Или создайте нового.
							</p>
							<button
								type="button"
								className={ui.primary}
								onClick={() => {
									setNewHero(true);
									setTab("hero");
								}}
							>
								Создать героя
							</button>
						</div>
					)}
				</div>
			</div>
			{pendingNavigation && (
				<NavigationConfirm
					onCancel={() => setPendingNavigation(null)}
					onDiscard={() => {
						setPendingNavigation(null);
						setDirty(false);
						setNotice("");
						pendingNavigation();
					}}
				/>
			)}
		</main>
	);
}

function NavigationConfirm({
	onCancel,
	onDiscard,
}: {
	onCancel: () => void;
	onDiscard: () => void;
}) {
	const dialog = useRef<HTMLDialogElement>(null);
	useEffect(() => {
		const element = dialog.current;
		element?.showModal();
		return () => element?.close();
	}, []);
	return (
		<dialog
			ref={dialog}
			className={ui.dialog}
			aria-labelledby="admin-navigation-title"
			onCancel={(event) => {
				event.preventDefault();
				onCancel();
			}}
		>
			<h2 id="admin-navigation-title">Есть несохранённые изменения</h2>
			<p>Изменённые поля и выбранные файлы будут потеряны при переходе.</p>
			<div className={ui.actions}>
				<button type="button" className={ui.primary} onClick={onCancel}>
					Продолжить редактирование
				</button>
				<button type="button" className={ui.secondary} onClick={onDiscard}>
					Перейти без сохранения
				</button>
			</div>
		</dialog>
	);
}
