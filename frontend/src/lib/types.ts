export type RecipeRes = {
	name: string;
	description: string;
	id: number;
};
export type Recipes = {
	recipes: Array<RecipeRes>;
};

export type Response = {
	id: number;
};

export type IngredientRes = {
	id: number;
	ingredient_id: number;
	name: string;
	amount: number;
	specifier: string;
	current_amount: string;
};
export type Ingredients = {
	ingredients: Array<IngredientRes>;
};
