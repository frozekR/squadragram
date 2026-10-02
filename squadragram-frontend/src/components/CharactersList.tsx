import CharacterCard from "./CharacterCard";
import { useCharacters } from "#/hooks/useCharacters";

export const CharactersList = () => {
	const {
		data: characters = [],
		isLoading,
		isError,
		refetch,
	} = useCharacters();
	if (isLoading)
		return (
			<p className="content-status" aria-live="polite">
				Loading heroes…
			</p>
		);
	if (isError)
		return (
			<div className="content-status" role="alert">
				<p>Could not load heroes.</p>
				<button type="button" onClick={() => void refetch()}>
					Try again
				</button>
			</div>
		);
	if (characters.length === 0)
		return <p className="content-status">No heroes yet.</p>;
	return (
		<div className="character-grid">
			{characters.map((character) => (
				<CharacterCard key={character.uuid} character={character} />
			))}
		</div>
	);
};
