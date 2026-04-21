import { PUBLIC_SERVER_URL, PUBLIC_SERVER_PORT } from "$env/static/public";
import { type IngredientRes } from "$lib/types";
const SERVER_URL = PUBLIC_SERVER_PORT
	? `${PUBLIC_SERVER_URL}:${PUBLIC_SERVER_PORT}`
	: PUBLIC_SERVER_URL;

export const postIngredient = async (id: number, formData: FormData): Promise<Response> => {
	const response = await fetch(`${SERVER_URL}/api/ingredient`, {
		method: "POST",
		body: JSON.stringify({
			id: id,
			name: formData.get("ingredient")?.toString(),
			amount: Number(formData.get("amount")),
			specifier: formData.get("specifier")?.toString(),
			current_amount: "out",
		}),
	});
	return response
}

export const deleteIngredientDB = async (id: number, item: IngredientRes): boolean => {
	const response = await fetch(`${SERVER_URL}/api/ingredient`, {
		method: "DELETE",
		body: JSON.stringify({
			recipe_id: id,
			id: item.ingredient_id,
			name: item.name.toLowerCase(),
		}),
	});

	return response.ok
};

