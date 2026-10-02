import { useState } from "react";

export function DemoVideo({
	src,
	title,
	poster,
}: {
	src: string;
	title: string;
	poster?: string;
}) {
	const [failedSrc, setFailedSrc] = useState<string>();
	return (
		<div className="demo-video">
			{failedSrc === src ? (
				<div className="video-unavailable" aria-live="polite">
					<p>Demo video is unavailable.</p>
					<button type="button" onClick={() => setFailedSrc(undefined)}>
						Try again
					</button>
				</div>
			) : (
				// biome-ignore lint/a11y/useMediaCaption: The media API currently has no caption assets for gameplay demos.
				<video
					key={src}
					src={src}
					poster={poster}
					controls
					playsInline
					preload="metadata"
					aria-label={title}
					onError={() => setFailedSrc(src)}
				>
					Your browser does not support video playback.
				</video>
			)}
		</div>
	);
}
