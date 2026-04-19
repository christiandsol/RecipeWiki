import { PUBLIC_URL } from "$env/static/public";
import type { IngredientRes, Ingredients } from "$lib/types";

export function useIngredients(id: number) {
	let listIngredients: Array<IngredientRes & { priority?: string }> = $state([]);

	$effect(() => {
		const get = async () => {
			const response = await fetch(`${PUBLIC_URL}/ingredients`, {
				method: "POST",
				body: JSON.stringify({ id }),
			});
			const data: Ingredients = await response.json();
			if (data.ingredients != null) {
				listIngredients = data.ingredients.map((item) => ({
					id: item.id,
					ingredient_id: item.ingredient_id,
					name: item.name,
					amount: item.amount,
					specifier: item.specifier,
					priority: item.current_amount,
				}));
			}
		};
		get();
	});

	return {
		get listIngredients() { return listIngredients; },
		set listIngredients(v) { listIngredients = v; },
	};
}
