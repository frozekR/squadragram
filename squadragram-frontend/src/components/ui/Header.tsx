import { useState, useEffect, useRef } from "react";

export default function Header() {
	const [isDropdownOpen, setIsDropdownOpen] = useState(false);
	const dropdownRef = useRef<HTMLDivElement | null>(null);

	// Обработчик клика вне меню
	useEffect(() => {
		function handleClickOutside(event: MouseEvent) {
			// Если кликнули не по кнопке и не по самому меню — закрываем
			if (
				dropdownRef.current &&
				event.target instanceof Node &&
				!dropdownRef.current.contains(event.target)
			) {
				setIsDropdownOpen(false);
			}
		}

		// Добавляем слушатель события при монтировании
		document.addEventListener("mousedown", handleClickOutside);

		// Убираем слушатель при размонтировании
		return () => {
			document.removeEventListener("mousedown", handleClickOutside);
		};
	}, []);

	return (
		<header
			className="
        w-[calc(100%-1.5rem)] max-w-4xl mx-auto mt-3 sm:mt-6
        bg-blue-400/70 backdrop-blur-md 
        border border-white/30 shadow-[0_8px_32px_0_rgba(31,38,135,0.37)] 
        rounded-3xl sm:rounded-full py-3 px-4 sm:px-8
        text-white font-medium transition-all duration-300
      "
		>
			<nav className="flex flex-wrap justify-center items-center gap-x-3 gap-y-2 sm:gap-x-6">
				<a href="/" className="hover:text-blue-200 transition-colors">
					Home
				</a>
				<span className="hidden sm:inline opacity-50">|</span>

				<a href="/heroes/" className="hover:text-blue-200 transition-colors">
					Heroes
				</a>
				<span className="hidden sm:inline opacity-50">|</span>
        
				{/* --- DROPDOWN С АНИМАЦИЕЙ И ЗАКРЫТИЕМ ПО КЛИКУ ВНЕ --- */}
				{/* Привязываем ref к обертке */}
				<div className="relative" ref={dropdownRef}>
					<button
						onClick={() => setIsDropdownOpen(!isDropdownOpen)}
						className="flex items-center gap-1 hover:text-blue-200 transition-colors focus:outline-none"
					>
						Links
						<svg
							className={`w-3 h-3 transition-transform duration-300 ${isDropdownOpen ? "rotate-180" : ""}`}
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
						>
							<path
								strokeLinecap="round"
								strokeLinejoin="round"
								strokeWidth="2"
								d="M19 9l-7 7-7-7"
							/>
						</svg>
					</button>

					{/* Само меню. Теперь оно рендерится всегда, но скрыто стилями при !isDropdownOpen */}
					<div
						className={`
              absolute top-full left-1/2 -translate-x-1/2 max-sm:left-auto max-sm:right-0 max-sm:translate-x-0 mt-4
              w-30
              bg-blue-400/70 backdrop-blur-lg 
              border border-white/30 
              rounded-3xl shadow-lg 
              overflow-hidden flex flex-col py-1 z-50
              
              /* Базовые классы для анимации */
              transition-all duration-300 ease-out origin-top
              
              /* Состояние открыто/закрыто */
              ${
								isDropdownOpen
									? "opacity-100 scale-100 translate-y-0 pointer-events-auto"
									: "opacity-0 scale-95 -translate-y-2 pointer-events-none"
							}
            `}
					>
						<a
							href="https://t.me/gekishinsquadra"
							className="px-4 py-2 text-sm hover:bg-white/20 transition-colors"
						>
							Telegram
						</a>
						<a
							href="#"
							className="px-4 py-2 text-sm hover:bg-white/20 transition-colors"
						>
							Support
						</a>
					</div>
				</div>
				{/* --- КОНЕЦ DROPDOWN --- */}
			</nav>
		</header>
	);
}
