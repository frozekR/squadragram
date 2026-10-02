import { adminUI as ui } from "#/components/admin/styles";
import { useEffect, useState, type FormEvent } from "react";
import { adminErrorMessage, saveHero, saveSkill } from "#/api/admin";
import type {
	Character,
	CharacterRole,
	CreateCharacterDTO,
} from "#/shared/types/character";
import {
	skillTypeLabels,
	type Skill,
	type SkillType,
} from "#/shared/types/skill";

export const roleLabels: Record<CharacterRole, string> = {
	DAMAGE: "Урон",
	TANK: "Танк",
	TECHNICAL: "Технический",
	MELEE: "Ближний бой",
	RANGED: "Дальний бой",
};

interface CommonProps {
	token: string;
	onDirty: (dirty: boolean) => void;
}

export function HeroEditor({
	token,
	character,
	onSaved,
	onDirty,
}: CommonProps & {
	character?: Character;
	onSaved: (hero: Character) => void;
}) {
	const initial: CreateCharacterDTO = {
		name: character?.name || "",
		description: character?.description || "",
		role: character?.role || "DAMAGE",
	};
	const [fields, setFields] = useState(initial);
	const [busy, setBusy] = useState(false);
	const [error, setError] = useState("");
	const dirty = JSON.stringify(fields) !== JSON.stringify(initial);
	useEffect(() => {
		onDirty(dirty);
		return () => onDirty(false);
	}, [dirty, onDirty]);
	async function submit(event: FormEvent) {
		event.preventDefault();
		setBusy(true);
		setError("");
		try {
			const saved = await saveHero(
				token,
				{
					...fields,
					name: fields.name.trim(),
					description: fields.description.trim(),
				},
				character?.uuid,
			);
			setFields({
				name: saved.name,
				description: saved.description,
				role: saved.role,
			});
			onSaved(saved);
		} catch (error) {
			setError(adminErrorMessage(error));
		} finally {
			setBusy(false);
		}
	}
	return (
		<section className={ui.panel}>
			<div className={ui.sectionHeading}>
				<div>
					<p className={ui.eyebrow}>Карточка героя</p>
					<h2>{character ? "Основная информация" : "Новый герой"}</h2>
				</div>
				{dirty && <span className={ui.unsaved}>Не сохранено</span>}
			</div>
			<form className={ui.form} onSubmit={submit}>
				<fieldset disabled={busy}>
					<label>
						Имя героя
						<input
							name="hero-name"
							value={fields.name}
							maxLength={96}
							required
							onChange={(e) => setFields({ ...fields, name: e.target.value })}
							placeholder="Например, Super Saiyan Vegeta"
						/>
					</label>
					<label>
						Роль
						<select
							value={fields.role}
							onChange={(e) =>
								setFields({ ...fields, role: e.target.value as CharacterRole })
							}
						>
							{Object.entries(roleLabels).map(([role, label]) => (
								<option key={role} value={role}>
									{label}
								</option>
							))}
						</select>
					</label>
					<label>
						Описание
						<textarea
							name="hero-description"
							value={fields.description}
							required
							rows={5}
							onChange={(e) =>
								setFields({ ...fields, description: e.target.value })
							}
							placeholder="Особенности героя и его стиль игры"
						/>
					</label>
				</fieldset>
				{error && (
					<p className={ui.error} role="alert">
						{error}
					</p>
				)}
				<div className={ui.actions}>
					<button
						type="submit"
						className={ui.primary}
						disabled={
							busy ||
							!fields.name.trim() ||
							!fields.description.trim() ||
							(Boolean(character) && !dirty)
						}
					>
						{busy
							? "Сохраняем…"
							: character
								? "Сохранить изменения"
								: "Создать героя"}
					</button>
					<button
						type="button"
						className={ui.secondary}
						disabled={busy || !dirty}
						onClick={() => {
							setFields(initial);
							setError("");
						}}
					>
						Сбросить
					</button>
				</div>
			</form>
		</section>
	);
}

export function SkillEditor({
	token,
	skill,
	owner,
	characters,
	defaultOrder = 0,
	onSaved,
	onDirty,
}: CommonProps & {
	skill?: Skill;
	owner: Character;
	characters: Character[];
	defaultOrder?: number;
	onSaved: (skill: Skill, ownerUUID: string) => void;
}) {
	const initial = {
		name: skill?.name || "",
		description: skill?.description || "",
		type: skill?.type || ("SKILL" as SkillType),
		character_uuid: owner.uuid,
		sort_order: String(skill?.sort_order ?? defaultOrder),
	};
	const [fields, setFields] = useState(initial);
	const [busy, setBusy] = useState(false);
	const [error, setError] = useState("");
	const dirty = JSON.stringify(fields) !== JSON.stringify(initial);
	useEffect(() => {
		onDirty(dirty);
		return () => onDirty(false);
	}, [dirty, onDirty]);
	async function submit(event: FormEvent) {
		event.preventDefault();
		setBusy(true);
		setError("");
		try {
			const saved = await saveSkill(
				token,
				{
					...fields,
					sort_order: Number(fields.sort_order),
					name: fields.name.trim(),
					description: fields.description.trim(),
				},
				skill?.uuid,
			);
			setFields({
				name: saved.name,
				description: saved.description,
				type: saved.type,
				character_uuid: fields.character_uuid,
				sort_order: String(saved.sort_order),
			});
			onSaved(saved, fields.character_uuid);
		} catch (error) {
			setError(adminErrorMessage(error));
		} finally {
			setBusy(false);
		}
	}
	return (
		<section className={ui.panel}>
			<div className={ui.sectionHeading}>
				<div>
					<p className={ui.eyebrow}>Скилл героя</p>
					<h2>{skill ? "Редактирование скилла" : "Новый скилл"}</h2>
				</div>
				{dirty && <span className={ui.unsaved}>Не сохранено</span>}
			</div>
			<form className={ui.form} onSubmit={submit}>
				<fieldset disabled={busy}>
					<label>
						Название
						<input
							name="skill-name"
							value={fields.name}
							maxLength={96}
							required
							onChange={(e) => setFields({ ...fields, name: e.target.value })}
							placeholder="Например, Big Bang Attack"
						/>
					</label>
					<label>
						Порядок в категории
						<input
							type="number"
							min={0}
							max={2147483647}
							step={1}
							required
							value={fields.sort_order}
							onChange={(event) =>
								setFields({ ...fields, sort_order: event.target.value })
							}
						/>
						<small className={ui.help}>
							Меньшее число — левее. При одинаковых значениях сохраняется
							порядок создания. Категория определяется полем «Тип».
						</small>
					</label>
					<div className={ui.formRow}>
						<label>
							Тип
							<select
								value={fields.type}
								onChange={(e) =>
									setFields({ ...fields, type: e.target.value as SkillType })
								}
							>
								{Object.entries(skillTypeLabels).map(([type, label]) => (
									<option key={type} value={type}>
										{label}
									</option>
								))}
							</select>
						</label>
						<label>
							Герой
							<select
								value={fields.character_uuid}
								onChange={(e) =>
									setFields({ ...fields, character_uuid: e.target.value })
								}
							>
								{characters.map((hero) => (
									<option key={hero.uuid} value={hero.uuid}>
										{hero.name}
									</option>
								))}
							</select>
						</label>
					</div>
					<label>
						Описание
						<textarea
							name="skill-description"
							value={fields.description}
							required
							rows={5}
							onChange={(e) =>
								setFields({ ...fields, description: e.target.value })
							}
							placeholder="Что делает скилл и как его использовать"
						/>
					</label>
				</fieldset>
				{error && (
					<p className={ui.error} role="alert">
						{error}
					</p>
				)}
				<div className={ui.actions}>
					<button
						type="submit"
						className={ui.primary}
						disabled={
							busy ||
							!fields.name.trim() ||
							!fields.description.trim() ||
							(Boolean(skill) && !dirty)
						}
					>
						{busy
							? "Сохраняем…"
							: skill
								? "Сохранить изменения"
								: "Создать скилл"}
					</button>
					<button
						type="button"
						className={ui.secondary}
						disabled={busy || !dirty}
						onClick={() => {
							setFields(initial);
							setError("");
						}}
					>
						Сбросить
					</button>
				</div>
			</form>
		</section>
	);
}
