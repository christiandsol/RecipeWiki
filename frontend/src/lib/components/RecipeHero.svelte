<script lang="ts">
    let { recipe, heroPreviewUrl, SERVER_URL, updateRecipe } = $props();
    let fileInput: HTMLInputElement;
</script>

<div class="recipe-hero">
    {#if recipe.image_url}
        <img
            class="recipe-hero-img"
            src={heroPreviewUrl ?? `${SERVER_URL}/images/${recipe.image_url}`}
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
            onkeydown={(e) => e.key === "Enter" && fileInput.click()}
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
            <span class="hero-upload-hint">PNG, JPG, WEBP · max 20 MB</span>
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
    ></textarea>
    <textarea
        class="recipe-caption-desc-input"
        value={recipe.description}
        onblur={(e) => updateRecipe("description", e.currentTarget.value)}
        rows="2"
        placeholder="Add a description"
    ></textarea>
</div>
