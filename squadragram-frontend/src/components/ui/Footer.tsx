export default function Footer() {
	return (
		<footer
			className="
        w-[calc(100%-1.5rem)] max-w-4xl mx-auto mb-3 sm:mb-6
        bg-slate-900/60 backdrop-blur-md 
        border border-white/20 
        rounded-3xl py-5 px-4 sm:py-6 sm:px-6 md:px-8
        text-center text-white 
        shadow-lg
      "
		>
			<div className="flex flex-col items-center justify-center gap-3">
				{/* Название и версия */}
				<h2 className="text-lg font-bold tracking-wide">
					squadra.gram{" "}
					<span className="text-white/70 font-medium text-base">v0.01</span>
				</h2>

				{/* Ссылки */}
				<div className="flex flex-wrap justify-center items-center gap-x-6 gap-y-2 text-sm font-medium">
					<a href="#" className="hover:text-blue-300 transition-colors">
						Privacy Policy
					</a>
				</div>

				{/* Дисклеймер (Авторские права) */}
				<p className="text-xs text-white/50 max-w-2xl mx-auto leading-relaxed mt-1 px-1">
					squadra.gram is not affiliated with or endorsed by Bandai Namco
					Entertainment Inc. / GANBARION
				</p>
			</div>
		</footer>
	);
}
