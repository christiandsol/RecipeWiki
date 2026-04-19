<script lang="ts">
    import { PUBLIC_URL } from "$env/static/public";
    import type { IngredientRes } from "$lib/types";

    let {
        listIngredients = $bindable(),
        onDeleteIngredient,
    }: {
        listIngredients: Array<IngredientRes & { priority?: string }>;
        onDeleteIngredient: (e: MouseEvent, item: IngredientRes) => void;
    } = $props();

    let openDropdown: number | null = $state(null);

    const priorities = [
        { value: "high", label: "High", color: "#f87171", bg: "#fef2f2" },
        { value: "medium", label: "Medium", color: "#fb923c", bg: "#fff7ed" },
        { value: "low", label: "Low", color: "#4ade80", bg: "#f0fdf4" },
        { value: "out", label: "Out", color: "#94a3b8", bg: "#f1f5f9" },
    ];

    const getPriority = (value: string) =>
        priorities.find((p) => p.value === value) ?? priorities[3];

    const toggleDropdown = (e: MouseEvent, ingredientId: number) => {
        e.stopPropagation();
        openDropdown = openDropdown === ingredientId ? null : ingredientId;
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
</script>

<svelte:window onclick={() => (openDropdown = null)} />

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
                        <span class="ingredient-unit">{item.specifier}</span>
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
                            onclick={(e) => onDeleteIngredient(e, item)}
                        >
                            Remove
                        </button>
                    </span>
                </li>
            {/each}
        </ul>
    {/if}
</section>
