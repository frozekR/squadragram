import { adminUI as ui } from "#/components/admin/styles";
import { useState, type FormEvent } from "react";
import { adminErrorMessage, verifyAdmin } from "#/api/admin";

export function AdminLogin({ onLogin }: { onLogin: (token: string) => void }) {
	const [token, setToken] = useState("");
	const [busy, setBusy] = useState(false);
	const [error, setError] = useState("");
	async function submit(event: FormEvent) {
		event.preventDefault();
		setBusy(true);
		setError("");
		try {
			await verifyAdmin(token.trim());
			onLogin(token.trim());
			setToken("");
		} catch (error) {
			setError(adminErrorMessage(error));
		} finally {
			setBusy(false);
		}
	}
	return (
		<main className={ui.login}>
			<div className={ui.loginCard}>
				<span className={ui.mark} aria-hidden="true">
					S
				</span>
				<p className={ui.eyebrow}>Squadragram · управление</p>
				<h1>Вход в админку</h1>
				<p className="text-sm leading-relaxed text-slate-500">
					Герои, скиллы и медиа — в одном месте.
				</p>
				<form onSubmit={submit} className={`${ui.form} mt-7 mb-5`}>
					<label>
						Административный токен
						<input
							type="password"
							name="admin-token"
							autoComplete="off"
							value={token}
							onChange={(event) => setToken(event.target.value)}
							required
							minLength={32}
							placeholder="Введите токен доступа"
							disabled={busy}
						/>
					</label>
					{error && (
						<p className={ui.error} role="alert">
							{error}
						</p>
					)}
					<button type="submit" className={ui.primary} disabled={busy}>
						{busy ? "Проверяем доступ…" : "Войти"}
					</button>
				</form>
				<p className={ui.help}>
					Токен хранится только в памяти этой страницы. После обновления
					потребуется войти снова.
				</p>
			</div>
		</main>
	);
}
