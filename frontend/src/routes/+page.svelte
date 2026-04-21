<script lang="ts">
	import "./home.css";
	import { goto } from "$app/navigation";
	import { type Recipes, type RecipeRes, type Response } from "$lib/types";
	import {
		PUBLIC_SERVER_URL,
		PUBLIC_SERVER_PORT,
		PUBLIC_URL,
	} from "$env/static/public";

	// URL
	const SERVER_URL = PUBLIC_SERVER_PORT
		? `${PUBLIC_SERVER_URL}:${PUBLIC_SERVER_PORT}`
		: PUBLIC_SERVER_URL;

	//Data Model
	let listRecipes: Array<{
		id: number;
		name: string | undefined;
		description: string | undefined;
		image_url: string | undefined;
	}> = $state([]);

	// Image upload state
	let previewUrl: string | null = $state(null);
	let isDragging = $state(false);
	let fileInput: HTMLInputElement;
	const FALLBACK =
		"https://images.unsplash.com/photo-1495521821757-a1efb6729352?w=400&q=80";
	const imageUrl = (filename: string | undefined) =>
		filename ? `${SERVER_URL}/images/${filename}` : FALLBACK;

	// Get recipes from server on mount
	$effect(() => {
		const getRecipes = async () => {
			const response = await fetch(`${SERVER_URL}/api/recipes`, {
				method: "GET",
			});
			const data: Recipes = await response.json();
			listRecipes = data.recipes.map((item: RecipeRes) => ({
				id: item.id,
				name: item.name,
				description: item.description,
				image_url: item.image_url,
			}));
		};
		getRecipes();
	});

	const onFileChange = async (img: File | null | undefined) => {
		if (!img) return;
		previewUrl = URL.createObjectURL(img);
	};

	const onDrop = (e: DragEvent) => {
		e.preventDefault();
		isDragging = false;
		const file = e.dataTransfer?.files[0];
		if (file && file.type.startsWith("image/")) {
			fileInput.files = e.dataTransfer!.files;
			onFileChange(file);
		}
	};

	const clearImage = () => {
		previewUrl = null;
		fileInput.value = "";
	};

	export const addRecipe = async (e: SubmitEvent) => {
		e.preventDefault();
		const formData = new FormData(e.target as HTMLFormElement);
		const response = await fetch(`${SERVER_URL}/api/recipe`, {
			method: "POST",
			body: formData,
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
				image_url: data.image_url,
			});
			clearImage();
			goto(`/recipe/${data.id}`);
		}
	};

	export const deleteRecipe = async (e: MouseEvent, recipeId: number) => {
		const response = await fetch(`${PUBLIC_URL}/api/recipe`, {
			method: "DELETE",
			body: JSON.stringify({ recipe_id: recipeId }),
		});
		const msg = await response.text();
		if (response.ok) {
			listRecipes = listRecipes.filter((r) => r.id !== recipeId);
		} else {
			console.log(
				`[ERROR] Error deleting recipe, server responded ${msg}`,
			);
		}
	};
</script>

<main class="home-page">
	<header class="site-header">
		<div class="site-header-top">
			<div>
				<h1 class="site-title">Sol's Recipes</h1>
				<p class="site-subtitle">Add and manage your recipes</p>
			</div>
			<a href="/fridge" class="fridge-btn">
				<svg
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="1.5"
				>
					<path
						d="M8 2h8a2 2 0 012 2v16a2 2 0 01-2 2H8a2 2 0 01-2-2V4a2 2 0 012-2z"
						stroke-linecap="round"
						stroke-linejoin="round"
					/>
					<path d="M6 10h12" stroke-linecap="round" />
					<path d="M10 6v2M10 14v3" stroke-linecap="round" />
				</svg>
				My Fridge
			</a>
		</div>
	</header>
	<section class="gallery-section">
		<h2 class="section-title">Your Recipes</h2>
		{#if listRecipes.length === 0}
			<p class="empty-state">No recipes yet. Add one below.</p>
		{:else}
			<div class="recipe-grid">
				{#each listRecipes as recipe}
					<div
						class="recipe-card"
						onclick={() => goto(`/recipe/${recipe.id}`)}
						role="button"
						tabindex="0"
						onkeydown={(e) =>
							e.key === "Enter" && goto(`/recipe/${recipe.id}`)}
					>
						<button
							class="recipe-card-delete"
							onclick={(e) => {
								e.stopPropagation();
								deleteRecipe(e, recipe.id);
							}}
							aria-label="Delete recipe">✕</button
						>
						<div class="recipe-card-image">
							<img
								src={imageUrl(recipe.image_url)}
								alt={recipe.name}
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
					</div>
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

			<!-- Image upload -->
			<input
				bind:this={fileInput}
				type="file"
				name="image"
				accept="image/*"
				class="file-input-hidden"
				onchange={(e) =>
					onFileChange((e.target as HTMLInputElement).files?.[0])}
			/>

			<div class="field-group">
				<label>Photo</label>
				{#if previewUrl}
					<div class="upload-preview">
						<img
							src={previewUrl}
							alt="Preview"
							class="preview-img"
						/>
						<button
							type="button"
							class="preview-clear"
							onclick={clearImage}
							aria-label="Remove image">✕</button
						>
					</div>
				{:else}
					<div
						class="upload-zone"
						class:upload-zone--drag={isDragging}
						ondragover={(e) => {
							e.preventDefault();
							isDragging = true;
						}}
						ondragleave={() => (isDragging = false)}
						ondrop={onDrop}
						onclick={() => fileInput.click()}
						role="button"
						tabindex="0"
						onkeydown={(e) =>
							e.key === "Enter" && fileInput.click()}
					>
						<svg
							class="upload-icon"
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
						<span class="upload-label"
							>Drop an image or <span class="upload-link"
								>browse</span
							></span
						>
						<span class="upload-hint"
							>PNG, JPG, WEBP · max 5 MB</span
						>
					</div>
				{/if}
			</div>

			<button class="submit-btn" type="submit">Add Recipe</button>
		</form>
	</section>
</main>
