const roleStyles: Record<
	string,
	{
		border: string;
		shadow: string;
	}
> = {
	ranged: {
		border: "border-blue-500",
		shadow: "shadow-[0_0_15px_rgba(59,130,246,0.5)]",
	},
	melee: {
		border: "border-red-500",
		shadow: "shadow-[0_0_15px_rgba(239,68,68,0.5)]",
	},
	default: {
		border: "border-gray-500",
		shadow: "shadow-md",
	},
};

const roleIcons: Record<string, React.ReactNode> = {
	melee: (
		<svg
			viewBox="0 0 24 24"
			fill="none"
			stroke="currentColor"
			strokeWidth="1.8"
			aria-hidden="true"
		>
			<path
				strokeLinecap="round"
				strokeLinejoin="round"
				d="M7 11V6a2 2 0 0 1 4 0v4V5a2 2 0 0 1 4 0v5V6a2 2 0 0 1 4 0v7.5A6.5 6.5 0 0 1 12.5 20H11a6 6 0 0 1-6-6v-3a2 2 0 0 1 4 0v1Zm0 0V8a2 2 0 0 1 4 0v3m4-1V8m4 2v-1"
			/>
		</svg>
	),
	ranged: (
		<svg
			viewBox="0 0 24 24"
			fill="none"
			stroke="currentColor"
			strokeWidth="1.8"
			aria-hidden="true"
		>
			<path
				strokeLinecap="round"
				strokeLinejoin="round"
				d="M12 3v18M3 12h18M5.6 5.6l12.8 12.8M18.4 5.6 5.6 18.4"
			/>
		</svg>
	),
};

interface CharacterCardProps {
	id: string | number;
	name: string;
	imageUrl: string;
	role?: string;
}

export default function CharacterCard({
	id,
	name,
	imageUrl,
	role = "default",
}: CharacterCardProps) {
	// Получаем стили в зависимости от роли (приводим к нижнему регистру на случай опечаток)
	const normalizedRole = role.toLowerCase();
	const styles = roleStyles[normalizedRole] || roleStyles.default;
	const roleIcon = roleIcons[normalizedRole];
	const linkHref = `/character/${id}`;

	return (
		<a
			href={linkHref}
			className={`
          relative block w-48 h-48 rounded-2xl overflow-hidden 
          border-4 transition-all duration-300 group
          ${styles.border} ${styles.shadow}
        `}
		>
			{/* Изображение персонажа */}
			<img
				src={imageUrl}
				alt={name}
				className="w-full h-full object-cover group-hover:scale-110 transition-transform duration-500 ease-out"
			/>

			{roleIcon && (
				<div className="absolute left-2 top-2 z-10 flex h-9 w-9 items-center justify-center rounded-full border-2 border-white/90 bg-black/45 p-1.5 text-white shadow-md backdrop-blur-sm">
					{roleIcon}
				</div>
			)}

			{/* Градиентная подложка снизу для читаемости текста */}
			<div className="absolute inset-0 bg-gradient-to-t from-black/90 via-black/20 to-transparent" />

			{/* Имя персонажа */}
			<div className="absolute bottom-3 left-0 right-0 text-center">
				<span className="text-white font-bold text-lg drop-shadow-md tracking-wide">
					{name}
				</span>
			</div>
		</a>
	);
}
