import { API_BASE_URL } from "./client";
import type { Character, CreateCharacterDTO } from "#/shared/types/character";
import type { Skill, SkillType } from "#/shared/types/skill";
import type { MediaMetadata, MediaType, OwnerType } from "#/shared/types/media";

export class AdminAPIError extends Error {
	constructor(
		message: string,
		public status: number,
	) {
		super(message);
	}
}

async function request<T>(
	path: string,
	token: string,
	options: RequestInit = {},
): Promise<T> {
	const headers = new Headers(options.headers);
	headers.set("Authorization", `Bearer ${token}`);
	if (typeof options.body === "string")
		headers.set("Content-Type", "application/json");
	const response = await fetch(`${API_BASE_URL}${path}`, {
		...options,
		headers,
	});
	if (!response.ok) {
		const body = await response.json().catch(() => null);
		throw new AdminAPIError(
			body?.error || `Ошибка запроса (${response.status})`,
			response.status,
		);
	}
	if (response.status === 204) return undefined as T;
	return response.json() as Promise<T>;
}

export function verifyAdmin(token: string) {
	return request<void>("/admin/verify", token, { method: "POST" });
}
export function saveHero(
	token: string,
	dto: CreateCharacterDTO,
	uuid?: string,
) {
	return request<Character>(
		uuid ? `/characters/${uuid}` : "/characters",
		token,
		{ method: uuid ? "PUT" : "POST", body: JSON.stringify(dto) },
	);
}
export interface SkillDTO {
	name: string;
	description: string;
	type: SkillType;
	character_uuid: string;
	sort_order?: number;
}
export function saveSkill(token: string, dto: SkillDTO, uuid?: string) {
	return request<Skill>(uuid ? `/skills/${uuid}` : "/skills", token, {
		method: uuid ? "PUT" : "POST",
		body: JSON.stringify(dto),
	});
}
export function uploadMedia(
	token: string,
	ownerType: OwnerType,
	ownerUUID: string,
	mediaType: MediaType,
	file: File,
	existingID?: number,
) {
	const body = new FormData();
	body.set("owner_type", ownerType);
	body.set("owner_uuid", ownerUUID);
	body.set("media_type", mediaType);
	body.set("file", file);
	return request<MediaMetadata>(
		existingID ? `/media-metadata/${existingID}` : "/media-metadata",
		token,
		{ method: existingID ? "PUT" : "POST", body },
	);
}
export function deleteMedia(token: string, id: number) {
	return request<void>(`/media-metadata/${id}`, token, { method: "DELETE" });
}

export function adminErrorMessage(error: unknown): string {
	if (error instanceof AdminAPIError) {
		if (error.status === 401)
			return "Доступ отклонён. Проверьте токен или войдите заново.";
		if (error.status === 409)
			return "Такая запись уже существует. Проверьте название или обновите данные.";
		if (error.status === 404) return "Запись не найдена. Обновите список.";
		if (error.status === 413) return "Файл превышает ограничение 32 МБ.";
		if (error.status === 415) return "Формат файла не поддерживается.";
		if (error.status === 502)
			return "Хранилище медиа недоступно. Повторите загрузку позже.";
		if (error.status === 503)
			return "Запись отключена: административный токен не настроен на сервере.";
		return error.message;
	}
	return "Не удалось связаться с сервером. Проверьте подключение и повторите попытку.";
}
