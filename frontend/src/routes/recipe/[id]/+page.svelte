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
		type Response,
		type IngredientRes,
		type GetStepResult,
	} from "$lib/types";
	import { PUBLIC_URL } from "$env/static/public";
	import IngredientList from "$lib/components/IngredientList.svelte";
	import StepList from "$lib/components/StepList.svelte";
	const SERVER_URL = PUBLIC_SERVER_PORT
		? `${PUBLIC_SERVER_URL}:${PUBLIC_SERVER_PORT}`
		: PUBLIC_SERVER_URL;
	const FALLBACK =
		"https://images.unsplash.com/photo-1495521821757-a1efb6729352?w=400&q=80";

	let { data: props } = $props();
	let listSteps: Array<{ id: number; text: string; stepNumber: number }> =
		$state([]);
	let heroPreviewUrl: string | null = $state(null);
	let draftText: string = $state("");
	let addingStep: boolean = $state(false);
	let editingStepId: number | null = $state(null);
	let editingText: string = $state("");

	console.log(`id passed: ${props.id}`);
	let listIngredients: Array<IngredientRes & { priority?: string }> = $state(
		[],
	);
	let openDropdown: string | null = $state(null);
	// get recipe image
	let recipe = $state(null);

	let draggedStepId: number | null = $state(null);
	let dragOverStepId: number | null = $state(null);

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

	$effect(() => {
		const getSteps = async () => {
			const response = await fetch(`${SERVER_URL}/steps/${props.id}`);
			if (!response.ok) return;
			const data = await response.json();
			if (data.steps != null) {
				listSteps = data.steps.map(
					(s: { step_id: number; step_text: string }) => ({
						id: s.step_id,
						text: s.step_text,
						stepNumber: s.step_number,
					}),
				);
			}
		};
		getSteps();
	});

	const updateImage = async (e: SubmitEvent) => {
		var formData = new FormData(e.target as HTMLFormElement);
		const response = await fetch(`${SERVER_URL}/recipe`, {
			method: "PATCH",
			body: formData,
		});
	};

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
				specifier: formData.get("specifier"),
				priority: "none",
			});
			(e.target as HTMLFormElement).reset();
		}
	};

	const deleteIngredient = async (e: MouseEvent, item: IngredientRes) => {
		e.preventDefault();
		const ok = await deleteIngredientDB(props.id, item);
		if (ok) {
			listIngredients = listIngredients.filter(
				(i) => i.ingredient_id !== item.ingredient_id,
			);
		} else {
			console.log("[ERROR] unable to delete ingredient");
		}
	};

	const setPriority = async (i: IngredientRes, priority: string) => {
		const response = await fetch(`${PUBLIC_URL}/ingredient`, {
			method: "PATCH",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify({ ...i, current_amount: priority }),
		});
		if (!response.ok) {
			console.log(`[ERROR] Unable to update priority`);
			return;
		}
		listIngredients = listIngredients.map((item) =>
			item.ingredient_id === i.ingredient_id
				? { ...item, priority, current_amount: priority }
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

	const handleAddStep = async () => {
		const data: GetStepResult = await addStep(props.id, draftText);
		if (!data) return;

		listSteps.push({
			id: data.step_id,
			text: draftText.trim(),
			stepNumber: data.step_number,
		});
		draftText = "";
		addingStep = false;
	};

	const handleUpdateStep = async (step_id: number, text: string) => {
		const ok: boolean = await updateStep(step_id, text);
		if (ok) {
			listSteps = listSteps.map((s) =>
				s.id === step_id ? { ...s, text: text.trim() } : s,
			);
			editingStepId = null;
		}
	};

	const deleteStep = async (stepId: number) => {
		const response = await fetch(`${SERVER_URL}/step`, {
			method: "DELETE",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify({ id: stepId }),
		});
		if (response.ok) {
			listSteps = listSteps.filter((s) => s.id !== stepId);
		}
	};
</script>

<svelte:window onclick={closeAllDropdowns} />

<main class="recipe-page">
	<SideNav />
	<div class="recipe-content" id="top">
		{#if recipe}
			<RecipeHero
				{recipe}
				{heroPreviewUrl}
				{SERVER_URL}
				{updateRecipe}
				onFileInputBind={(el) => (fileInput = el)}
			/>
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
