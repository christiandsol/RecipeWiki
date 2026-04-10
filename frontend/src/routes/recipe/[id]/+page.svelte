<script lang="ts">
	import "./recipe.css";
	import { Ingredient } from "$lib/ingredients";
	type IngredientRes = {
		id?: number;
		name: string;
		amount: number;
		specifier: number;
	};
	type Ingredients = {
		ingredients: Array<IngredientRes>;
	};
	let { data: props } = $props();
	console.log(`RECEIVED PROP: data: ${props.id}`);
	const PORT = 8080;
	const URL = `http://localhost:${PORT}`;
	let listIngredients: Array<{
		name: string;
		amount: string;
		specifier: string;
	}> = $state([]);
	$effect(() => {
		const send = async () => {
			const response = await fetch(`${URL}/ingredients`, {
				method: "POST",
				body: JSON.stringify({
					id: props.id,
				}),
			});
			const ingredients: Ingredients = await response.json();
			listIngredients = ingredients.ingredients.map((item) => ({
				name: item.name,
				amount: String(item.amount),
				specifier: item.specifier,
			}));
		};
		send();
	});
	const addIngredient = async (e: SubmitEvent) => {
		e.preventDefault();
		const formData = new FormData(e.target as HTMLFormElement);
		const response = await fetch(`${URL}/ingredient`, {
			method: "POST",
			body: JSON.stringify({
				id: props.id,
				name: formData.get("ingredient")?.toString().toLowerCase(),
				amount: Number(formData.get("amount")),
				specifier: formData.get("specifier")?.toString().toLowerCase(),
			}),
		});
		if (response.ok) {
			listIngredients.push({
				name: formData.get("ingredient") as string,
				amount: formData.get("amount") as string,
				specifier: formData.get("specifier") as string,
			});
		}
	};
</script>

<main class="page">
	<section class="form-section">
		<h2 class="section-title">Add Ingredient</h2>
		<form class="ingredient-form" onsubmit={addIngredient}>
			<div class="field-group">
				<label for="ingredient">Name</label>
				<input
					id="ingredient"
					type="text"
					name="ingredient"
					placeholder="e.g. flour"
					required
				/>
			</div>
			<div class="field-row">
				<div class="field-group">
					<label for="amount">Amount</label>
					<input
						id="amount"
						type="number"
						name="amount"
						placeholder="0"
						required
					/>
				</div>
				<div class="field-group">
					<label for="specifier">Unit</label>
					<select id="specifier" name="specifier">
						<optgroup label="Volume">
							<option value="tsp">Teaspoon (tsp)</option>
							<option value="tbsp">Tablespoon (tbsp)</option>
							<option value="cups">Cups</option>
							<option value="ml">Milliliters (ml)</option>
							<option value="l">Liters (l)</option>
						</optgroup>
						<optgroup label="Weight">
							<option value="grams">Grams (g)</option>
							<option value="kg">Kilograms (kg)</option>
							<option value="oz">Ounces (oz)</option>
							<option value="lb">Pounds (lb)</option>
						</optgroup>
						<optgroup label="Count">
							<option value="unit">Unit / Piece</option>
							<option value="serving">Serving</option>
						</optgroup>
					</select>
				</div>
			</div>
			<button class="submit-btn" type="submit">Add Ingredient</button>
		</form>
	</section>

	<section class="list-section">
		<h2 class="section-title">Ingredients</h2>
		{#if listIngredients.length === 0}
			<p class="empty-state">No ingredients added yet.</p>
		{:else}
			<ul class="ingredient-list">
				{#each listIngredients as item}
					<li class="ingredient-item">
						<span class="ingredient-name">{item.name}</span>
						<span class="ingredient-measure">
							<span class="ingredient-amount">{item.amount}</span>
							<span class="ingredient-unit">{item.specifier}</span
							>
						</span>
					</li>
				{/each}
			</ul>
		{/if}
	</section>
</main>
