import { PUBLIC_SERVER_URL, PUBLIC_SERVER_PORT } from "$env/static/public";

const SERVER_URL = PUBLIC_SERVER_PORT
	? `${PUBLIC_SERVER_URL}:${PUBLIC_SERVER_PORT}`
	: PUBLIC_SERVER_URL;

export const useRecipe = (id: number) => {
	let recipe = $state(null)
	let heroPreviewUrl: string | null = $state(null);
	$effect(() => {
		const getRecipe = async () => {
			const response = await fetch(`${SERVER_URL}/api/recipe/${id}`);
			const recipe = await response.json();
		};
		getRecipe();
	});
	return {
		get recipe() { return recipe; },
		set recipe(v) { recipe = v; },
		get heroPreviewUrl() { return heroPreviewUrl; },
		set heroPreviewUrl(v) { heroPreviewUrl = v; },
	};
}
