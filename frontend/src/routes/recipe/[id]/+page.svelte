<script lang="ts">
	import "./recipe.css";
	import { PUBLIC_SERVER_URL, PUBLIC_SERVER_PORT } from "$env/static/public";
	import {
		type Ingredients,
		type Response,
		type IngredientRes,
	} from "$lib/types";
	import { PUBLIC_URL } from "$env/static/public";
	const SERVER_URL = `${PUBLIC_SERVER_URL}:${PUBLIC_SERVER_PORT}`;
	const FALLBACK =
		"https://images.unsplash.com/photo-1495521821757-a1efb6729352?w=400&q=80";

	let { data: props } = $props();
	console.log(`id passed: ${props.id}`);
	let listIngredients: Array<IngredientRes & { priority?: string }> = $state(
		[],
	);
	let openDropdown: string | null = $state(null);
	// get recipe image
	let recipe = $state(null);

	$effect(() => {
		const getRecipe = async () => {
			const response = await fetch(`${SERVER_URL}/recipe/${props.id}`, {
				method: "GET",
			});
			const data = await response.json();
			recipe = data;
		};
		getRecipe();
	});

	// get ingredients
	$effect(() => {
		const send = async () => {
			const response = await fetch(`${PUBLIC_URL}/ingredients`, {
				method: "POST",
				body: JSON.stringify({ id: props.id }),
			});
			const ingredients: Ingredients = await response.json();
			if (ingredients.ingredients != null) {
				listIngredients = ingredients.ingredients.map((item) => ({
					id: item.id,
					ingredient_id: item.ingredient_id,
					name: item.name,
					amount: item.amount,
					specifier: item.specifier,
					priority: item.current_amount,
				}));
			}
		};
		send();
	});

	const addIngredient = async (e: SubmitEvent) => {
		e.preventDefault();
		const formData = new FormData(e.target as HTMLFormElement);
		const response = await fetch(`${PUBLIC_URL}/ingredient`, {
			method: "POST",
			body: JSON.stringify({
				id: props.id,
				name: formData.get("ingredient")?.toString(),
				amount: Number(formData.get("amount")),
				specifier: formData.get("specifier")?.toString(),
				current_amount: "out",
			}),
		});

		const data: Response = await response.json();
		if (response.ok) {
			listIngredients.push({
				id: data.id,
				name: formData.get("ingredient") as string,
				amount: Number(formData.get("amount")),
				specifier: formData.get("specifier"),
				priority: "none",
			});
		}
	};

	const deleteIngredient = async (e: MouseEvent, item: IngredientRes) => {
		e.preventDefault();
		const response = await fetch(`${PUBLIC_URL}/ingredient`, {
			method: "DELETE",
			body: JSON.stringify({
				recipe_id: props.id,
				id: item.ingredient_id,
				name: item.name.toLowerCase(),
			}),
		});

		const msg = await response.text();
		if (response.ok) {
			listIngredients = listIngredients.filter(
				(i) => i.ingredient_id !== item.ingredient_id,
			);
		} else {
			console.log(
				`[ERROR] Error deleting ingredient, server responded ${msg}`,
			);
		}
	};

	const setPriority = async (i: IngredientRes, priority: string) => {
		i.current_amount = priority;
		console.log(`Sending ingredient with amount: ${i.current_amount}`);
		const response = await fetch(`${PUBLIC_URL}/ingredient`, {
			method: "PATCH",
			headers: {
				"Content-Type": "application/json",
			},
			body: JSON.stringify(i),
		});
		const msg = await response.text();
		if (!response.ok) {
			console.log(`Successfully updated Ingredient`);
			console.log(`[ERROR] Unable to update amount: ${msg}`);
			return;
		}
		listIngredients = listIngredients.map((item) =>
			item.ingredient_id === i.ingredient_id
				? { ...item, priority }
				: item,
		);
		openDropdown = null;
	};

	const toggleDropdown = (e: MouseEvent, ingredientId: string) => {
		e.stopPropagation();
		openDropdown = openDropdown === ingredientId ? null : ingredientId;
	};

	const closeAllDropdowns = () => {
		openDropdown = null;
	};

	const priorities = [
		{ value: "high", label: "High", color: "#f87171", bg: "#fef2f2" },
		{ value: "medium", label: "Medium", color: "#fb923c", bg: "#fff7ed" },
		{ value: "low", label: "Low", color: "#4ade80", bg: "#f0fdf4" },
		{ value: "none", label: "None", color: "#94a3b8", bg: "#f1f5f9" },
	];

	const getPriority = (value: string) =>
		priorities.find((p) => p.value === value) ?? priorities[3];

	const updateRecipe = async (
		field: "name" | "description",
		value: string,
	) => {
		console.log("Sending PATCH to", `${SERVER_URL}/recipe`);
		const response = await fetch(`${SERVER_URL}/recipe`, {
			method: "PATCH",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify({
				id: recipe.id,
				name: recipe.name,
				description: recipe.description,
				[field]: value, // override whichever field changed
			}),
		});

		if (!response.ok) {
			console.error("[ERROR] Failed to update recipe");
			return;
		}

		recipe = { ...recipe, [field]: value };
	};
</script>

<svelte:window onclick={closeAllDropdowns} />

<main class="page">
	{#if recipe}
		<div class="recipe-hero">
			<img
				class="recipe-hero-img"
				src={recipe.image_url
					? `${SERVER_URL}/images/${recipe.image_url}`
					: FALLBACK}
				alt={recipe.name}
			/>
		</div>
		<div class="recipe-caption">
			<textarea
				class="recipe-caption-name-input"
				value={recipe.name}
				onblur={(e) => updateRecipe("name", e.target.value)}
				rows="1"
			/>
			<!-- swap this for an input later when you add edit mode -->
			<textarea
				class="recipe-caption-desc-input"
				value={recipe.description}
				onblur={(e) => updateRecipe("description", e.target.value)}
				rows="2"
			/>
		</div>
	{/if}
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
					{@const p = getPriority(item.priority ?? "none")}
					<li class="ingredient-item">
						<span class="ingredient-name">{item.name}</span>
						<span class="ingredient-measure">
							<span class="ingredient-amount">{item.amount}</span>
							<span class="ingredient-unit">{item.specifier}</span
							>

							<!-- Priority dropdown -->
							<div class="priority-wrapper">
								<button
									class="priority-badge"
									style="color: {p.color}; background: {p.bg};"
									onclick={(e) =>
										toggleDropdown(e, item.ingredient_id)}
								>
									<span
										class="priority-dot"
										style="background:{p.color};"
									></span>
									{p.label}
									<svg
										class="priority-chevron"
										viewBox="0 0 10 6"
										fill="none"
									>
										<path
											d="M1 1l4 4 4-4"
											stroke="currentColor"
											stroke-width="1.5"
											stroke-linecap="round"
										/>
									</svg>
								</button>

								{#if openDropdown === item.ingredient_id}
									<ul
										class="priority-menu"
										onclick={(e) => e.stopPropagation()}
									>
										{#each priorities as opt}
											<li>
												<button
													class="priority-option"
													class:priority-option--active={item.priority ===
														opt.value}
													style="--opt-color:{opt.color}; --opt-bg:{opt.bg};"
													onclick={() =>
														setPriority(
															item,
															opt.value,
														)}
												>
													<span
														class="priority-dot"
														style="background:{opt.color};"
													></span>
													{opt.label}
												</button>
											</li>
										{/each}
									</ul>
								{/if}
							</div>

							<button
								class="remove-btn"
								onclick={(e) => deleteIngredient(e, item)}
								>Remove</button
							>
						</span>
					</li>
				{/each}
			</ul>
		{/if}
	</section>
</main>
