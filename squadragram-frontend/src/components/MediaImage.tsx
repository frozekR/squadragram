import { useState } from "react";

interface MediaImageProps {
	src?: string;
	alt: string;
	className?: string;
	fallback?: string;
	fit?: "cover" | "contain";
	loading?: "eager" | "lazy";
}

export function MediaImage({
	src,
	alt,
	className = "",
	fallback = "No image",
	fit = "cover",
	loading = "lazy",
}: MediaImageProps) {
	const [failedSrc, setFailedSrc] = useState<string>();
	return (
		<div className={`media-image ${className}`}>
			{src && src !== failedSrc ? (
				<img
					src={src}
					alt={alt}
					loading={loading}
					className={fit === "contain" ? "media-contain" : "media-cover"}
					onError={() => setFailedSrc(src)}
				/>
			) : (
				<div
					className="media-placeholder"
					role="img"
					aria-label={`${alt}: image unavailable`}
				>
					<span>{fallback}</span>
				</div>
			)}
		</div>
	);
}
