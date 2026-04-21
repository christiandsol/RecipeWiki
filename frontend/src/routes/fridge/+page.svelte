<script lang="ts">
	import "./fridge.css";
	import {
		PUBLIC_SERVER_URL,
		PUBLIC_SERVER_PORT,
		PUBLIC_URL,
	} from "$env/static/public";
	import type { Ingredients, IngredientRes } from "$lib/types";

	const SERVER_URL = PUBLIC_SERVER_PORT
		? `${PUBLIC_SERVER_URL}:${PUBLIC_SERVER_PORT}`
		: PUBLIC_SERVER_URL;

	let listIngredients: Array<IngredientRes & { priority?: string }> = $state(
		[],
	);
	let loading = $state(true);

	const priorities = [
		{ value: "high", label: "High", color: "#f87171", bg: "#fef2f2" },
		{ value: "medium", label: "Medium", color: "#fb923c", bg: "#fff7ed" },
		{ value: "low", label: "Low", color: "#4ade80", bg: "#f0fdf4" },
		{ value: "out", label: "Out", color: "#94a3b8", bg: "#f1f5f9" },
	];

	const getPriority = (value: string) =>
		priorities.find((p) => p.value === value) ?? priorities[3];

	let openDropdown: number | null = $state(null);

	const toggleDropdown = (e: MouseEvent, id: number) => {
		e.stopPropagation();
		openDropdown = openDropdown === id ? null : id;
	};

	$effect(() => {
		const fetchAll = async () => {
			const response = await fetch(`${SERVER_URL}/api/fridge`);
			if (!response.ok) {
				loading = false;
				return;
			}
			const data: Ingredients = await response.json();
			if (data.ingredients != null) {
				listIngredients = data.ingredients.map((item) => ({
					...item,
					priority: item.current_amount,
				}));
			}
			loading = false;
		};
		fetchAll();
	});

	const setPriority = async (i: IngredientRes, priority: string) => {
		const response = await fetch(`${PUBLIC_URL}/api/ingredient`, {
			method: "PATCH",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify({ ...i, current_amount: priority }),
		});
		if (!response.ok) return;
		listIngredients = listIngredients.map((item) =>
			item.ingredient_id === i.ingredient_id
				? { ...item, priority, current_amount: priority }
				: item,
		);
		openDropdown = null;
	};

	const grouped = $derived(
		listIngredients.reduce(
			(acc, item) => {
				const key = item.priority ?? "out";
				if (!acc[key]) acc[key] = [];
				acc[key].push(item);
				return acc;
			},
			{} as Record<string, typeof listIngredients>,
		),
	);
</script>

<svelte:window onclick={() => (openDropdown = null)} />

<main class="fridge-page">
	<header class="fridge-header">
		<a href="/" class="back-btn">
			<svg
				viewBox="0 0 16 16"
				fill="none"
				stroke="currentColor"
				stroke-width="1.8"
			>
				<path
					d="M10 3L5 8l5 5"
					stroke-linecap="round"
					stroke-linejoin="round"
				/>
			</svg>
			Back
		</a>
		<div>
			<h1 class="fridge-title">My Fridge</h1>
			<p class="fridge-subtitle">Track what you have and what you need</p>
		</div>
	</header>

	{#if loading}
		<p class="fridge-loading">Loading…</p>
	{:else if listIngredients.length === 0}
		<div class="fridge-empty">
			<p class="fridge-empty-label">Your fridge is empty</p>
			<p class="fridge-empty-hint">
				Add ingredients to your recipes to see them here
			</p>
		</div>
	{:else}
		<div class="fridge-groups">
			{#each priorities as group}
				{#if grouped[group.value]?.length}
					<section class="fridge-group">
						<div class="fridge-group-header">
							<span
								class="fridge-group-dot"
								style="background:{group.color};"
							></span>
							<h2 class="fridge-group-title">{group.label}</h2>
							<span class="fridge-group-count"
								>{grouped[group.value].length}</span
							>
						</div>
						<ul class="fridge-list">
							{#each grouped[group.value] as item}
								{@const p = getPriority(item.priority ?? "out")}
								<li class="fridge-item">
									<div class="fridge-item-info">
										<span class="fridge-item-name"
											>{item.name}</span
										>
										<span class="fridge-item-amount"
											>{item.amount}
											{item.specifier}</span
										>
									</div>
									<div class="priority-wrapper">
										<button
											class="priority-badge"
											style="color:{p.color}; background:{p.bg};"
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
												onclick={(e) =>
													e.stopPropagation()}
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
								</li>
							{/each}
						</ul>
					</section>
				{/if}
			{/each}
		</div>
	{/if}
</main>
