<script lang="ts">
	import "./recipe.css";
	import { PUBLIC_SERVER_URL, PUBLIC_SERVER_PORT } from "$env/static/public";
	import { addStep, updateStep, reorderStep } from "$lib/api/stepService.js";
	import { patchRecipe, fetchRecipe } from "$lib/api/recipeService.js";
	import {
		deleteIngredientDB,
		postIngredient,
	} from "$lib/api/ingredientService.js";
	import SideNav from "$lib/components/SideNav.svelte";
	import RecipeHero from "$lib/components/RecipeHero.svelte";
	import IngredientForm from "$lib/components/IngredientForm.svelte";

	import {
		type Ingredients,
		type IngredientRes,
		type RecipeRes,
	} from "$lib/types";
	import { PUBLIC_URL } from "$env/static/public";
	import IngredientList from "$lib/components/IngredientList.svelte";
	import StepList from "$lib/components/StepList.svelte";
	const SERVER_URL = PUBLIC_SERVER_PORT
		? `${PUBLIC_SERVER_URL}:${PUBLIC_SERVER_PORT}`
		: PUBLIC_SERVER_URL;

	let { data: props } = $props();
	let listSteps: Array<{ id: number; text: string; stepNumber: number }> =
		$state([]);
	let heroPreviewUrl: string | null = $state(null);

	let listIngredients: Array<IngredientRes> = $state([]);
	// get recipe image
	let recipe: RecipeRes | null = $state(null);

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
					current_amount: item.current_amount,
					specifier: item.specifier,
					priority: item.current_amount,
				}));
			}
		};
		send();
	});

	$effect(() => {
		const getSteps = async () => {
			const response = await fetch(`${SERVER_URL}/steps/${props.id}`);
			if (!response.ok) return;
			const data = await response.json();
			if (data.steps != null) {
				listSteps = data.steps.map(
					(s: {
						step_id: number;
						step_text: string;
						step_number: number;
					}) => ({
						id: s.step_id,
						text: s.step_text,
						stepNumber: s.step_number,
					}),
				);
			}
		};
		getSteps();
	});

	const addIngredient = async (e: SubmitEvent) => {
		console.log(`HERE`);
		e.preventDefault();
		const formData = new FormData(e.target as HTMLFormElement);
		const response = await postIngredient(props.id, formData);
		const ingredient = await response.json();
		if (response.ok) {
			listIngredients.push({
				id: props.id,
				ingredient_id: ingredient.id,
				name: formData.get("ingredient") as string,
				amount: Number(formData.get("amount")),
				specifier: String(formData.get("specifier")),
				current_amount: "none",
			});
			(e.target as HTMLFormElement).reset();
		}
	};

	const deleteIngredient = async (e: MouseEvent, item: IngredientRes) => {
		e.preventDefault();
		const ok = deleteIngredientDB(props.id, item);
		if (ok) {
			listIngredients = listIngredients.filter(
				(i) => i.ingredient_id !== item.ingredient_id,
			);
		} else {
			console.log("[ERROR] unable to delete ingredient");
		}
	};

	const updateRecipe = async (
		field: "name" | "description" | "image",
		value: string | File,
	) => {
		if (recipe == null) return;
		if (field === "image") {
			heroPreviewUrl = URL.createObjectURL(value as File);
		}

		const ok = await patchRecipe(
			recipe.id,
			recipe.name,
			recipe.description,
			field,
			value,
		);

		if (!ok) {
			console.error("[ERROR] Failed to update recipe");
			heroPreviewUrl = null;
			return;
		}

		if (field === "image") {
			recipe = await fetchRecipe(recipe.id);
			heroPreviewUrl = null;
		} else {
			recipe = { ...recipe, [field]: value };
		}
	};

	const deleteStep = async (stepId: number): Promise<boolean> => {
		const response = await fetch(`${SERVER_URL}/step`, {
			method: "DELETE",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify({ id: stepId }),
		});
		return response.ok;
	};
</script>

<main class="recipe-page">
	<SideNav />
	<div class="recipe-content" id="top">
		{#if recipe}
			<RecipeHero {recipe} {heroPreviewUrl} {SERVER_URL} {updateRecipe} />
		{/if}

		<IngredientForm {addIngredient} />
		<IngredientList
			bind:listIngredients
			onDeleteIngredient={deleteIngredient}
		/>
		<StepList
			bind:listSteps
			onAddStep={(text) => addStep(props.id, text)}
			onUpdateStep={updateStep}
			onDeleteStep={deleteStep}
			onReorderStep={reorderStep}
		/>
	</div>
</main>
