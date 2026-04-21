import { PUBLIC_SERVER_URL, PUBLIC_SERVER_PORT } from "$env/static/public";
const SERVER_URL = PUBLIC_SERVER_PORT
	? `${PUBLIC_SERVER_URL}:${PUBLIC_SERVER_PORT}`
	: PUBLIC_SERVER_URL;

export const patchRecipe = async (
	recipeId: number,
	recipeName: string,
	recipeDescription: string,
	field: "name" | "description" | "image",
	value: string | File,
): Promise<boolean> => {
	const formData = new FormData();
	formData.append("id", String(recipeId));
	formData.append("name", field === "name" ? (value as string) : recipeName);
	formData.append("description", field === "description" ? (value as string) : recipeDescription);
	if (field === "image") {
		formData.append("image", value as File);
	}
	const response = await fetch(`${SERVER_URL}/api/recipe`, {
		method: "PATCH",
		body: formData,
	});
	return response.ok;
};

export const fetchRecipe = async (recipeId: number) => {
	const response = await fetch(`${SERVER_URL}/api/recipe/${recipeId}`);
	if (!response.ok) return null;
	return response.json();
};
