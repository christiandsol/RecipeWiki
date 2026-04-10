<script lang="ts">
	import "./home.css";

	type RecipeRes = {
		name: string;
		description: string;
		id: number;
	};
	type Recipes = {
		recipes: Array<RecipeRes>;
	};
	type Response = {
		id: number;
	};
	import { goto } from "$app/navigation";
	import { get } from "svelte/store";
	let listRecipes: Array<{
		id: number;
		name: string | undefined;
		description: string | undefined;
	}> = $state([]);
	const PORT = 8080;
	const URL = `http://localhost:${PORT}`;
	$effect(() => {
		const getRecipes = async () => {
			console.log("FETCHING");
			const response = await fetch(`${URL}/recipes`, { method: "GET" });
			const data: Recipes = await response.json();
			listRecipes = data.recipes.map((item: RecipeRes) => ({
				id: item.id,
				name: item.name,
				description: item.description,
			}));
		};
		getRecipes();
	});
	const addRecipe = async (e: SubmitEvent) => {
		e.preventDefault();
		const formData = new FormData(e.target as HTMLFormElement);
		const response = await fetch(`${URL}/recipe`, {
			method: "POST",
			body: JSON.stringify({
				name: formData.get("name")?.toString().toLowerCase(),
				description: formData
					.get("description")
					?.toString()
					.toLowerCase(),
			}),
		});
		if (!response.ok) {
			console.error(
				`[ERROR] Failed to add recipe with status ${response.status}`,
			);
		} else {
			const data: Response = await response.json();
			listRecipes.push({
				id: data.id,
				name: formData.get("name")?.toString().toLowerCase(),
				description: formData
					.get("description")
					?.toString()
					.toLowerCase(),
			});
			goto(`/recipe/${data.id}`);
		}
	};
</script>

<main class="page">
	<header class="site-header">
		<h1 class="site-title">Sol's Recipes</h1>
		<p class="site-subtitle">Add and manage your recipes</p>
	</header>

	<section class="gallery-section">
		<h2 class="section-title">Your Recipes</h2>
		{#if listRecipes.length === 0}
			<p class="empty-state">No recipes yet. Add one below.</p>
		{:else}
			<div class="recipe-grid">
				{#each listRecipes as recipe}
					<button
						class="recipe-card"
						onclick={() => goto(`/recipe/${recipe.id}`)}
					>
						<div class="recipe-card-image">
							<img
								src="https://images.unsplash.com/photo-1495521821757-a1efb6729352?w=400&q=80"
								alt="Recipe"
							/>
						</div>
						<div class="recipe-card-body">
							<h3 class="recipe-card-name">{recipe.name}</h3>
							{#if recipe.description}
								<p class="recipe-card-description">
									{recipe.description}
								</p>
							{/if}
						</div>
					</button>
				{/each}
			</div>
		{/if}
	</section>

	<section class="form-section">
		<h2 class="section-title">New Recipe</h2>
		<form class="recipe-form" onsubmit={addRecipe}>
			<div class="field-group">
				<label for="name">Recipe Name</label>
				<input
					id="name"
					type="text"
					name="name"
					placeholder="e.g. banana bread"
					required
				/>
			</div>
			<div class="field-group">
				<label for="description">Description</label>
				<input
					id="description"
					type="text"
					name="description"
					placeholder="e.g. a moist and sweet loaf"
				/>
			</div>
			<button class="submit-btn" type="submit">Add Recipe</button>
		</form>
	</section>
</main>
