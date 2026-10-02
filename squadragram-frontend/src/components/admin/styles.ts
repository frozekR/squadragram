// Static Tailwind utilities shared by admin components. No global admin CSS.
export const adminUI = {
	workspace:
		"mx-auto mt-6 mb-12 w-[calc(100%-1.5rem)] max-w-[1360px] min-w-0 flex-1 text-slate-900 sm:mt-9 sm:mb-16 sm:w-[calc(100%-2.5rem)]",
	topbar:
		"mb-6 flex flex-wrap items-center justify-between gap-4 [&_h1]:text-3xl [&_h1]:font-extrabold [&_h1]:tracking-tight sm:mb-7",
	access:
		"inline-flex items-center gap-1.5 text-xs font-semibold text-emerald-700 [&_span]:text-[9px]",
	layout:
		"grid min-w-0 grid-cols-1 items-start gap-5 lg:grid-cols-[270px_minmax(0,1fr)] lg:gap-6 xl:grid-cols-[310px_minmax(0,1fr)]",
	sidebar:
		"min-w-0 rounded-2xl border border-slate-200 bg-white p-4 lg:sticky lg:top-6 lg:max-h-[calc(100dvh-3rem)] lg:overflow-y-auto lg:p-5",
	sidebarHeading:
		"flex flex-wrap items-center justify-between gap-3 [&_h2]:text-lg [&_h2]:font-bold [&_h2_span]:ml-1 [&_h2_span]:text-xs [&_h2_span]:text-slate-400",
	main: "grid min-w-0 gap-5",
	panel:
		"min-w-0 rounded-2xl border border-slate-200 bg-white p-4 [&_h2]:text-xl [&_h2]:font-bold [&_h2]:leading-snug [&_h2]:wrap-anywhere sm:p-6 xl:p-7",
	entityHeading:
		"flex min-w-0 flex-wrap items-center justify-between gap-3 py-1 [&_h2]:text-xl [&_h2]:font-bold [&_h2]:wrap-anywhere [&_p]:mt-1.5 [&_p]:text-sm [&_p]:text-slate-500 sm:[&_h2]:text-2xl",
	actions: "flex flex-wrap items-center gap-2.5 max-sm:[&>button]:flex-1",
	primary:
		"inline-flex min-h-11 items-center justify-center gap-1.5 rounded-xl border px-4 py-2.5 text-sm font-semibold leading-snug cursor-pointer transition-colors focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-sky-500 disabled:cursor-not-allowed disabled:opacity-50 border-sky-600 bg-sky-600 text-white enabled:hover:bg-sky-700",
	secondary:
		"inline-flex min-h-11 items-center justify-center gap-1.5 rounded-xl border px-4 py-2.5 text-sm font-semibold leading-snug cursor-pointer transition-colors focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-sky-500 disabled:cursor-not-allowed disabled:opacity-50 border-slate-200 bg-slate-50 text-slate-600 enabled:hover:bg-slate-100",
	danger:
		"inline-flex min-h-11 items-center justify-center gap-1.5 rounded-xl border px-4 py-2.5 text-sm font-semibold leading-snug cursor-pointer transition-colors focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-sky-500 disabled:cursor-not-allowed disabled:opacity-50 border-rose-200 bg-rose-50 text-rose-700 enabled:hover:bg-rose-100",
	search:
		"mt-5 mb-4 block text-xs font-semibold text-slate-500 [&_input]:mt-2 [&_input]:block [&_input]:min-h-11 [&_input]:w-full [&_input]:min-w-0 [&_input]:rounded-xl [&_input]:border [&_input]:border-slate-300 [&_input]:bg-slate-50 [&_input]:px-3 [&_input]:py-2.5 [&_input]:text-base [&_input]:font-normal [&_input]:text-slate-900 [&_input]:outline-none [&_input]:focus:border-sky-500 [&_input]:focus:ring-2 [&_input]:focus:ring-sky-100 [&_input]:disabled:opacity-60 [&_input]:sm:text-sm",
	heroList:
		"grid max-h-52 gap-1.5 overflow-y-auto overscroll-contain lg:max-h-none lg:overflow-visible",
	heroItem:
		"flex min-h-14 w-full min-w-0 cursor-pointer items-center gap-3 rounded-xl border border-transparent p-2.5 text-left hover:bg-slate-50 aria-pressed:border-sky-200 aria-pressed:bg-sky-50 focus-visible:outline-2 focus-visible:outline-sky-500 [&>span]:min-w-0 [&_strong]:block [&_strong]:text-sm [&_strong]:font-semibold [&_strong]:wrap-anywhere [&_small]:mt-1 [&_small]:block [&_small]:text-xs [&_small]:text-slate-500",
	avatar: "h-11 w-11 shrink-0 rounded-xl [&_.media-placeholder]:p-1!",
	sectionHeading:
		"mb-5 flex min-w-0 flex-wrap items-center justify-between gap-3 max-sm:[&>button]:w-full",
	form: "grid min-w-0 gap-5 [&_fieldset]:grid [&_fieldset]:min-w-0 [&_fieldset]:gap-5 [&_fieldset]:border-0 [&_fieldset]:p-0 [&_label]:block [&_label]:text-sm [&_label]:font-semibold [&_label]:text-slate-600 [&_input]:mt-2 [&_input]:block [&_input]:min-h-11 [&_input]:w-full [&_input]:min-w-0 [&_input]:rounded-xl [&_input]:border [&_input]:border-slate-300 [&_input]:bg-slate-50 [&_input]:px-3 [&_input]:py-2.5 [&_input]:text-base [&_input]:font-normal [&_input]:text-slate-900 [&_input]:outline-none [&_input]:focus:border-sky-500 [&_input]:focus:ring-2 [&_input]:focus:ring-sky-100 [&_input]:disabled:opacity-60 [&_input]:sm:text-sm [&_select]:mt-2 [&_select]:block [&_select]:min-h-11 [&_select]:w-full [&_select]:min-w-0 [&_select]:rounded-xl [&_select]:border [&_select]:border-slate-300 [&_select]:bg-slate-50 [&_select]:px-3 [&_select]:py-2.5 [&_select]:text-base [&_select]:font-normal [&_select]:text-slate-900 [&_select]:outline-none [&_select]:focus:border-sky-500 [&_select]:focus:ring-2 [&_select]:focus:ring-sky-100 [&_select]:disabled:opacity-60 [&_select]:sm:text-sm [&_textarea]:mt-2 [&_textarea]:block [&_textarea]:min-h-11 [&_textarea]:w-full [&_textarea]:min-w-0 [&_textarea]:rounded-xl [&_textarea]:border [&_textarea]:border-slate-300 [&_textarea]:bg-slate-50 [&_textarea]:px-3 [&_textarea]:py-2.5 [&_textarea]:text-base [&_textarea]:font-normal [&_textarea]:text-slate-900 [&_textarea]:outline-none [&_textarea]:focus:border-sky-500 [&_textarea]:focus:ring-2 [&_textarea]:focus:ring-sky-100 [&_textarea]:disabled:opacity-60 [&_textarea]:sm:text-sm [&_textarea]:min-h-32 [&_textarea]:resize-y [&_textarea]:leading-relaxed",
	formRow: "grid min-w-0 grid-cols-1 gap-4 sm:grid-cols-2",
	help: "text-xs leading-relaxed text-slate-500",
	error:
		"rounded-xl bg-rose-50 px-3 py-3 text-sm leading-relaxed text-rose-700 wrap-anywhere [&_button]:ml-2 [&_button]:min-h-11 [&_button]:cursor-pointer [&_button]:underline",
	success: "text-xs text-emerald-700",
	unsaved:
		"rounded-full bg-amber-50 px-2.5 py-1 text-xs whitespace-nowrap text-amber-700",
	banner:
		"flex items-start justify-between gap-3 rounded-xl border border-emerald-200 bg-emerald-50 p-3 text-sm text-emerald-700 wrap-anywhere sm:p-4 [&_button]:flex [&_button]:size-11 [&_button]:shrink-0 [&_button]:cursor-pointer [&_button]:items-center [&_button]:justify-center [&_button]:text-xl",
	tabs: "flex min-w-0 gap-1 border-b border-slate-200 [&_button]:min-h-11 [&_button]:flex-1 [&_button]:cursor-pointer [&_button]:border-b-2 [&_button]:border-transparent [&_button]:px-3 [&_button]:py-3 [&_button]:text-sm [&_button]:text-slate-500 [&_button]:focus-visible:outline-2 [&_button]:focus-visible:outline-sky-500 [&_button[aria-pressed=true]]:border-sky-600 [&_button[aria-pressed=true]]:font-semibold [&_button[aria-pressed=true]]:text-sky-600 sm:[&_button]:flex-none sm:[&_button]:px-5",
	skillList:
		"grid min-w-0 gap-4 [&_button]:flex [&_button]:min-h-14 [&_button]:w-full [&_button]:min-w-0 [&_button]:cursor-pointer [&_button]:items-center [&_button]:gap-3 [&_button]:rounded-xl [&_button]:border [&_button]:border-slate-200 [&_button]:p-3 [&_button]:text-left [&_button]:focus-visible:outline-2 [&_button]:focus-visible:outline-sky-500 [&_button[aria-pressed=true]]:border-sky-300 [&_button[aria-pressed=true]]:bg-sky-50 [&_button>span]:min-w-0 [&_button>span:last-child]:ml-auto [&_button>span:last-child]:shrink-0 [&_button>span:last-child]:text-slate-400 [&_strong]:block [&_strong]:text-sm [&_strong]:font-semibold [&_strong]:wrap-anywhere [&_small]:mt-1 [&_small]:block [&_small]:text-xs [&_small]:text-slate-500",
	skillCategory: "grid min-w-0 gap-2.5",
	mediaGrid:
		"mt-5 grid min-w-0 grid-cols-1 gap-4 min-[600px]:grid-cols-2 xl:grid-cols-3",
	mediaSlot:
		"grid min-w-0 content-start gap-3 rounded-2xl border border-slate-200 bg-slate-50 p-4",
	mediaTitle:
		"flex flex-wrap items-center justify-between gap-2 [&_h3]:text-sm [&_h3]:font-bold [&_span]:text-xs [&_span]:text-slate-500",
	mediaPreview:
		"h-48 overflow-hidden rounded-xl bg-slate-100 [&>.media-image]:h-full [&_.demo-video]:h-full [&_video]:h-full! [&_video]:max-h-full! sm:h-40",
	emptyMedia: "grid h-full place-items-center text-sm text-slate-500",
	fileLabel:
		"block min-w-0 text-sm font-semibold text-sky-700 [&_input]:mt-2 [&_input]:block [&_input]:min-h-11 [&_input]:w-full [&_input]:min-w-0 [&_input]:text-xs [&_input]:font-normal [&_input]:text-slate-500 [&_input]:file:mr-2 [&_input]:file:min-h-11 [&_input]:file:cursor-pointer [&_input]:file:rounded-lg [&_input]:file:border [&_input]:file:border-slate-300 [&_input]:file:bg-white [&_input]:file:px-3 [&_input]:file:text-xs [&_input]:file:text-slate-600 [&_input]:focus-visible:outline-2 [&_input]:focus-visible:outline-sky-500",
	fileInfo: "text-xs text-slate-500 wrap-anywhere",
	confirm: "grid gap-3 border-t border-rose-200 pt-3 text-sm text-rose-700",
	login: "mx-auto my-8 w-[calc(100%-1.5rem)] max-w-lg flex-1 sm:my-16",
	loginCard:
		"rounded-2xl border border-slate-200 bg-white p-5 shadow-sm [&_h1]:mb-3 [&_h1]:text-2xl [&_h1]:font-extrabold [&_h1]:tracking-tight [&_h1]:text-slate-900 sm:p-9",
	mark: "mb-6 inline-grid size-13 place-items-center rounded-2xl bg-sky-600 text-3xl font-extrabold text-white",
	welcome:
		"py-10 text-center sm:py-16 [&>p:last-of-type]:mx-auto [&>p:last-of-type]:mt-4 [&>p:last-of-type]:mb-6 [&>p:last-of-type]:max-w-md [&>p:last-of-type]:text-sm [&>p:last-of-type]:leading-relaxed [&>p:last-of-type]:text-slate-500",
	dialog:
		"m-auto max-h-[calc(100dvh-2rem)] w-[calc(100%-2rem)] max-w-lg overflow-y-auto rounded-2xl border border-slate-200 bg-white p-5 text-slate-900 shadow-2xl backdrop:bg-slate-900/40 backdrop:backdrop-blur-sm [&_h2]:mb-3 [&_h2]:text-xl [&_h2]:font-bold [&_p]:mb-6 [&_p]:text-sm [&_p]:leading-relaxed [&_p]:text-slate-500 sm:p-7",
	eyebrow:
		"mb-2 text-[11px] font-extrabold uppercase tracking-widest text-slate-500",
	status: "py-6 text-center text-sm leading-relaxed text-slate-500",
} as const;
