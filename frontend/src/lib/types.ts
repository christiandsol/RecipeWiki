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
	name: string;
	amount: number;
	specifier: number;
};
export type Ingredients = {
	ingredients: Array<IngredientRes>;
};

export type Response = {
	id: number;
};
