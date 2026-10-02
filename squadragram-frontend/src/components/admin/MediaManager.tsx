import { adminUI as ui } from "#/components/admin/styles";
import { useEffect, useRef, useState } from "react";
import { adminErrorMessage, deleteMedia, uploadMedia } from "#/api/admin";
import { mediaFileURL } from "#/api/media";
import { DemoVideo } from "#/components/DemoVideo";
import { MediaImage } from "#/components/MediaImage";
import type { MediaMetadata, MediaType, OwnerType } from "#/shared/types/media";

const mediaLabels: Record<MediaType, string> = {
	ICON: "Иконка",
	RENDER: "Рендер",
	DEMO: "Демо-видео",
};
interface Props {
	token: string;
	ownerType: OwnerType;
	ownerUUID: string;
	metadata: MediaMetadata[];
	onChanged: () => Promise<void>;
	onDirty: (key: string, dirty: boolean) => void;
}

export function MediaManager(props: Props) {
	return (
		<section className={ui.panel}>
			<div className={ui.sectionHeading}>
				<div>
					<p className={ui.eyebrow}>Изображения и видео</p>
					<h2>Медиа</h2>
				</div>
			</div>
			<p className={ui.help}>
				По одному файлу каждого типа. Новый файл заменит текущий после успешной
				загрузки.
			</p>
			<div className={ui.mediaGrid}>
				{(["ICON", "RENDER", "DEMO"] as MediaType[]).map((type) => (
					<MediaSlot
						key={`${props.ownerUUID}-${type}`}
						{...props}
						type={type}
						existing={props.metadata.find((media) => media.media_type === type)}
					/>
				))}
			</div>
		</section>
	);
}

function MediaSlot({
	token,
	ownerType,
	ownerUUID,
	onChanged,
	onDirty,
	type,
	existing,
}: Props & { type: MediaType; existing?: MediaMetadata }) {
	const [file, setFile] = useState<File>();
	const [preview, setPreview] = useState("");
	const [busy, setBusy] = useState(false);
	const [error, setError] = useState("");
	const [notice, setNotice] = useState("");
	const [confirmRemove, setConfirmRemove] = useState(false);
	const input = useRef<HTMLInputElement>(null);
	useEffect(() => {
		const key = `${ownerUUID}-${type}`;
		onDirty(key, Boolean(file) || busy);
		return () => onDirty(key, false);
	}, [file, busy, ownerUUID, type, onDirty]);
	useEffect(() => {
		if (!file) {
			setPreview("");
			return;
		}
		const url = URL.createObjectURL(file);
		setPreview(url);
		return () => URL.revokeObjectURL(url);
	}, [file]);
	const src = preview || (existing ? mediaFileURL(existing) : undefined);
	async function upload() {
		if (!file) return;
		setBusy(true);
		setError("");
		setNotice("");
		try {
			await uploadMedia(token, ownerType, ownerUUID, type, file, existing?.id);
			await onChanged();
			setFile(undefined);
			if (input.current) input.current.value = "";
			setNotice("Файл сохранён");
		} catch (error) {
			setError(adminErrorMessage(error));
		} finally {
			setBusy(false);
		}
	}
	async function remove() {
		if (!existing) return;
		setBusy(true);
		setError("");
		setNotice("");
		try {
			await deleteMedia(token, existing.id);
			await onChanged();
			setConfirmRemove(false);
			setFile(undefined);
			if (input.current) input.current.value = "";
			setNotice("Файл удалён");
		} catch (error) {
			setError(adminErrorMessage(error));
		} finally {
			setBusy(false);
		}
	}
	return (
		<article className={ui.mediaSlot}>
			<div className={ui.mediaTitle}>
				<h3>{mediaLabels[type]}</h3>
				<span>
					{file ? "Предпросмотр" : existing ? "Загружено" : "Нет файла"}
				</span>
			</div>
			<div className={ui.mediaPreview}>
				{type === "DEMO" ? (
					src ? (
						<DemoVideo src={src} title="Предпросмотр демо" />
					) : (
						<div className={ui.emptyMedia}>Видео не загружено</div>
					)
				) : (
					<MediaImage
						src={src}
						alt={mediaLabels[type]}
						fit="contain"
						fallback="Изображение не загружено"
					/>
				)}
			</div>
			<label className={ui.fileLabel}>
				{existing ? "Выбрать новый файл" : "Выбрать файл"}
				<input
					ref={input}
					type="file"
					disabled={busy}
					accept={type === "DEMO" ? ".mp4,.webm" : ".png,.jpg,.jpeg,.gif,.webp"}
					onChange={(event) => {
						const next = event.target.files?.[0];
						setError("");
						setNotice("");
						setConfirmRemove(false);
						if (next && (next.size > 32 * 1024 * 1024 || next.size === 0)) {
							setFile(undefined);
							setError(
								next.size === 0 ? "Файл пуст." : "Максимальный размер — 32 МБ.",
							);
							event.target.value = "";
							return;
						}
						setFile(next);
					}}
				/>
			</label>
			<p className={ui.fileInfo}>
				{file
					? `${file.name} · ${(file.size / 1024 / 1024).toFixed(2)} МБ`
					: type === "DEMO"
						? "MP4 / WebM · до 32 МБ"
						: "PNG / JPEG / GIF / WebP · до 32 МБ"}
			</p>
			{error && (
				<p className={ui.error} role="alert">
					{error}
				</p>
			)}
			{notice && (
				<p className={ui.success} aria-live="polite">
					{notice}
				</p>
			)}
			<div className={ui.actions}>
				<button
					type="button"
					className={ui.primary}
					onClick={() => void upload()}
					disabled={!file || busy}
				>
					{busy ? "Обрабатываем…" : existing ? "Заменить" : "Загрузить"}
				</button>
				{file && (
					<button
						type="button"
						className={ui.secondary}
						disabled={busy}
						onClick={() => {
							setFile(undefined);
							if (input.current) input.current.value = "";
						}}
					>
						Отменить
					</button>
				)}
				{existing && !file && (
					<button
						type="button"
						className={ui.danger}
						disabled={busy}
						onClick={() => setConfirmRemove(true)}
					>
						Удалить
					</button>
				)}
			</div>
			{confirmRemove && (
				<div className={ui.confirm} role="alert">
					<p>
						Удалить {mediaLabels[type].toLowerCase()}? Файл будет удалён из
						хранилища.
					</p>
					<div className={ui.actions}>
						<button
							type="button"
							className={ui.danger}
							disabled={busy}
							onClick={() => void remove()}
						>
							Да, удалить
						</button>
						<button
							type="button"
							className={ui.secondary}
							disabled={busy}
							onClick={() => setConfirmRemove(false)}
						>
							Оставить
						</button>
					</div>
				</div>
			)}
		</article>
	);
}
