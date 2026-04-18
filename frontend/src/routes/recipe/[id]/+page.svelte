<script lang="ts">
	import "./recipe.css";
	import { PUBLIC_SERVER_URL, PUBLIC_SERVER_PORT } from "$env/static/public";
	import {
		type Ingredients,
		type Response,
		type IngredientRes,
	} from "$lib/types";
	import { PUBLIC_URL } from "$env/static/public";
	const SERVER_URL = PUBLIC_SERVER_PORT
		? `${PUBLIC_SERVER_URL}:${PUBLIC_SERVER_PORT}`
		: PUBLIC_SERVER_URL;
	const FALLBACK =
		"https://images.unsplash.com/photo-1495521821757-a1efb6729352?w=400&q=80";

	let { data: props } = $props();
	let listSteps: Array<{ id: number; text: string }> = $state([]);
	let heroPreviewUrl: string | null = $state(null);

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

	const updateImage = async (e: SubmitEvent) => {
		var formData = new FormData(e.target as HTMLFormElement);
		const response = await fetch(`${SERVER_URL}/recipe`, {
			method: "PATCH",
			body: formData,
		});
	};

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

	let fileInput: HTMLInputElement;

	const updateRecipe = async (
		field: "name" | "description" | "image",
		value: string | File,
	) => {
		const formData = new FormData();
		formData.append("id", String(recipe.id));
		formData.append(
			"name",
			field === "name" ? (value as string) : recipe.name,
		);
		formData.append(
			"description",
			field === "description" ? (value as string) : recipe.description,
		);

		if (field === "image") {
			heroPreviewUrl = URL.createObjectURL(value as File); // instant preview
			formData.append("image", value as File);
		}

		const response = await fetch(`${SERVER_URL}/recipe`, {
			method: "PATCH",
			body: formData,
		});

		if (!response.ok) {
			console.error("[ERROR] Failed to update recipe");
			heroPreviewUrl = null;
			return;
		}

		if (field === "image") {
			const updated = await fetch(`${SERVER_URL}/recipe/${recipe.id}`);
			recipe = await updated.json();
			heroPreviewUrl = null;
		} else {
			recipe = { ...recipe, [field]: value };
		}
	};

	const addStep = async (e: Event) => {
		console.log("ADDING STEP");
	};
</script>

<svelte:window onclick={closeAllDropdowns} />

<main class="recipe-page">
	<!-- Floating side nav -->
	<nav class="side-nav">
		<a href="#top" class="side-nav-item" title="Top">
			<svg
				viewBox="0 0 16 16"
				fill="none"
				stroke="currentColor"
				stroke-width="1.8"
			>
				<path
					d="M8 12V4M4 7l4-4 4 4"
					stroke-linecap="round"
					stroke-linejoin="round"
				/>
			</svg>
			<span class="side-nav-label">Top</span>
		</a>
		<div class="side-nav-divider"></div>
		<a href="#ingredients" class="side-nav-item" title="Ingredients">
			<svg
				viewBox="0 0 16 16"
				fill="none"
				stroke="currentColor"
				stroke-width="1.8"
			>
				<circle cx="8" cy="8" r="5.5" />
				<path d="M8 5.5v5M5.5 8h5" stroke-linecap="round" />
			</svg>
			<span class="side-nav-label">Ingredients</span>
		</a>
		<a href="#steps" class="side-nav-item" title="Steps">
			<svg
				viewBox="0 0 16 16"
				fill="none"
				stroke="currentColor"
				stroke-width="1.8"
			>
				<path d="M3 4h10M3 8h7M3 12h5" stroke-linecap="round" />
			</svg>
			<span class="side-nav-label">Steps</span>
		</a>
	</nav>

	<div class="recipe-content" id="top">
		{#if recipe}
			<div class="recipe-hero">
				{#if recipe.image_url}
					<img
						class="recipe-hero-img"
						src={heroPreviewUrl ??
							`${SERVER_URL}/images/${recipe.image_url}`}
						alt={recipe.name}
					/>
					<div class="recipe-hero-overlay">
						<button
							class="hero-change-btn"
							type="button"
							onclick={() => fileInput.click()}
						>
							Change photo
						</button>
					</div>
				{:else}
					<div
						class="recipe-hero-upload"
						role="button"
						tabindex="0"
						onclick={() => fileInput.click()}
						onkeydown={(e) =>
							e.key === "Enter" && fileInput.click()}
					>
						<div class="hero-upload-icon-wrap">
							<svg
								viewBox="0 0 24 24"
								fill="none"
								stroke="currentColor"
								stroke-width="1.5"
							>
								<path
									d="M3 16.5v2.25A2.25 2.25 0 005.25 21h13.5A2.25 2.25 0 0021 18.75V16.5m-13.5-9L12 3m0 0l4.5 4.5M12 3v13.5"
									stroke-linecap="round"
									stroke-linejoin="round"
								/>
							</svg>
						</div>
						<span class="hero-upload-label">Add a photo</span>
						<span class="hero-upload-hint"
							>PNG, JPG, WEBP · max 20 MB</span
						>
					</div>
				{/if}
			</div>

			<input
				bind:this={fileInput}
				type="file"
				accept="image/*"
				class="file-input-hidden"
				onchange={(e) => {
					const file = (e.target as HTMLInputElement).files?.[0];
					if (file) updateRecipe("image", file);
				}}
			/>

			<div class="recipe-caption">
				<textarea
					class="recipe-caption-name-input"
					value={recipe.name}
					onblur={(e) => updateRecipe("name", e.currentTarget.value)}
					rows="1"
					placeholder="Recipe name"
				/>
				<textarea
					class="recipe-caption-desc-input"
					value={recipe.description}
					onblur={(e) =>
						updateRecipe("description", e.currentTarget.value)}
					rows="2"
					placeholder="Add a description"
				/>
			</div>
		{/if}

		<!-- Add Ingredient -->
		<section class="form-section" id="ingredients">
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

		<!-- Ingredients list -->
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
								<span class="ingredient-amount"
									>{item.amount}</span
								>
								<span class="ingredient-unit"
									>{item.specifier}</span
								>
								<div class="priority-wrapper">
									<button
										class="priority-badge"
										style="color: {p.color}; background: {p.bg};"
										onclick={(e) =>
											toggleDropdown(
												e,
												item.ingredient_id,
											)}
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
								>
									Remove
								</button>
							</span>
						</li>
					{/each}
				</ul>
			{/if}
		</section>

		<!-- Steps -->
		<section class="steps-section" id="steps">
			<div class="steps-header">
				<h2 class="section-title">Steps</h2>
				<button
					class="add-step-btn"
					onclick={(e) => {
						addStep(e);
					}}
				>
					<svg
						viewBox="0 0 16 16"
						fill="none"
						stroke="currentColor"
						stroke-width="2"
					>
						<path d="M8 3v10M3 8h10" stroke-linecap="round" />
					</svg>
					Add step
				</button>
			</div>

			{#if listSteps.length === 0}
				<div class="steps-empty">
					<div class="steps-empty-icon">
						<svg
							viewBox="0 0 24 24"
							fill="none"
							stroke="currentColor"
							stroke-width="1.5"
						>
							<path
								d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"
								stroke-linecap="round"
								stroke-linejoin="round"
							/>
						</svg>
					</div>
					<p class="steps-empty-label">No steps yet</p>
					<p class="steps-empty-hint">
						Add step-by-step instructions for this recipe
					</p>
				</div>
			{:else}
				<ol class="steps-list">
					{#each listSteps as step, i}
						<li class="step-item">
							<div class="step-number">{i + 1}</div>
							<div class="step-body">
								<p class="step-text">{step.text}</p>
								<div class="step-actions">
									<button class="step-action-btn">Edit</button
									>
									<button
										class="step-action-btn step-action-btn--danger"
										>Remove</button
									>
								</div>
							</div>
						</li>
					{/each}
				</ol>
			{/if}
		</section>
	</div>
</main>
