<script lang="ts">
    let {
        listIngredients,
        getPriority,
        toggleDropdown,
        openDropdown,
        priorities,
        setPriority,
        deleteIngredient,
    } = $props();
</script>

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
