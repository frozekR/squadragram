import { API_BASE_URL } from "./client";
import type { MediaMetadata, MediaType } from "#/shared/types/media";

export function mediaFileURL(media: MediaMetadata): string {
	// MinIO object_key stays private; the backend serves files by metadata ID.
	const version = encodeURIComponent(media.object_key.split("/").pop() || "");
	return `${API_BASE_URL}/media-metadata/${media.id}/file?v=${version}`;
}

export function findMediaURL(
	metadata: MediaMetadata[] | null | undefined,
	type: MediaType,
): string | undefined {
	const media = metadata?.find((item) => item.media_type === type);
	return media ? mediaFileURL(media) : undefined;
}
