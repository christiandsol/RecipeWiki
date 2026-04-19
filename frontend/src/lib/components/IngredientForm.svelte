<script lang="ts">
    let { addIngredient } = $props();
</script>

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
