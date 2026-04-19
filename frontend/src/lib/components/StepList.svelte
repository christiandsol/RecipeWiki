<script lang="ts">
    let {
        listSteps = $bindable(),
        onAddStep,
        onUpdateStep,
        onDeleteStep,
        onReorderStep,
    }: {
        listSteps: Array<{ id: number; text: string; stepNumber: number }>;
        onAddStep: (
            text: string,
        ) => Promise<{ step_id: number; step_number: number } | null>;
        onUpdateStep: (id: number, text: string) => Promise<boolean>;
        onDeleteStep: (id: number) => Promise<boolean>;
        onReorderStep: (
            id: number,
            before: number,
            after: number,
        ) => Promise<void>;
    } = $props();

    let draftText: string = $state("");
    let addingStep: boolean = $state(false);
    let editingStepId: number | null = $state(null);
    let editingText: string = $state("");
    let draggedStepId: number | null = $state(null);
    let dragOverStepId: number | null = $state(null);

    const handleAddStep = async () => {
        if (!draftText.trim()) return;
        const text = draftText.trim();
        const data = await onAddStep(text);
        if (!data) return;
        listSteps = [
            ...listSteps,
            { id: data.step_id, text, stepNumber: data.step_number },
        ];
        draftText = "";
        addingStep = false;
    };
    const handleUpdateStep = async (id: number, text: string) => {
        const ok = await onUpdateStep(id, text);
        if (ok) {
            listSteps = listSteps.map((s) =>
                s.id === id ? { ...s, text: text.trim() } : s,
            );
            editingStepId = null;
        }
    };

    const handleDeleteStep = async (id: number) => {
        const ok = await onDeleteStep(id);
        if (ok) {
            listSteps = listSteps.filter((s) => s.id !== id);
        }
    };

    const handleDrop = (
        e: DragEvent,
        i: number,
        step: { id: number; stepNumber: number },
    ) => {
        e.preventDefault();
        if (draggedStepId === null || draggedStepId === step.id) return;

        const fromIndex = listSteps.findIndex((s) => s.id === draggedStepId);
        const toIndex = i;
        const before = toIndex === 0 ? 0 : listSteps[toIndex - 1].stepNumber;
        const after =
            toIndex >= listSteps.length - 1
                ? listSteps[listSteps.length - 1].stepNumber + 1
                : listSteps[toIndex].stepNumber;

        const updated = [...listSteps];
        const [moved] = updated.splice(fromIndex, 1);
        updated.splice(toIndex, 0, {
            ...moved,
            stepNumber: (before + after) / 2,
        });
        listSteps = updated;

        onReorderStep(draggedStepId, before, after);
        draggedStepId = null;
        dragOverStepId = null;
    };
</script>

<section class="steps-section" id="steps">
    <div class="steps-header">
        <h2 class="section-title">Steps</h2>
        <button
            class="add-step-btn"
            onclick={() => {
                addingStep = true;
            }}
            disabled={addingStep}
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

    {#if listSteps.length === 0 && !addingStep}
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
                <li
                    class="step-item"
                    class:step-item--dragging={draggedStepId === step.id}
                    class:step-item--dragover={dragOverStepId === step.id}
                    draggable="true"
                    ondragstart={() => {
                        draggedStepId = step.id;
                    }}
                    ondragend={() => {
                        draggedStepId = null;
                        dragOverStepId = null;
                    }}
                    ondragover={(e) => {
                        e.preventDefault();
                        if (draggedStepId !== step.id) dragOverStepId = step.id;
                    }}
                    ondragleave={() => {
                        if (dragOverStepId === step.id) dragOverStepId = null;
                    }}
                    ondrop={(e) => handleDrop(e, i, step)}
                >
                    <div class="step-drag-handle" aria-hidden="true">
                        <svg
                            viewBox="0 0 16 16"
                            fill="none"
                            stroke="currentColor"
                            stroke-width="1.5"
                        >
                            <circle
                                cx="5.5"
                                cy="4"
                                r="1"
                                fill="currentColor"
                                stroke="none"
                            />
                            <circle
                                cx="5.5"
                                cy="8"
                                r="1"
                                fill="currentColor"
                                stroke="none"
                            />
                            <circle
                                cx="5.5"
                                cy="12"
                                r="1"
                                fill="currentColor"
                                stroke="none"
                            />
                            <circle
                                cx="10.5"
                                cy="4"
                                r="1"
                                fill="currentColor"
                                stroke="none"
                            />
                            <circle
                                cx="10.5"
                                cy="8"
                                r="1"
                                fill="currentColor"
                                stroke="none"
                            />
                            <circle
                                cx="10.5"
                                cy="12"
                                r="1"
                                fill="currentColor"
                                stroke="none"
                            />
                        </svg>
                    </div>
                    <div class="step-number">{i + 1}</div>
                    <div class="step-body">
                        {#if editingStepId !== null && editingStepId === step.id}
                            <textarea
                                class="step-edit-input"
                                bind:value={editingText}
                                rows="3"
                                onkeydown={(e) => {
                                    if (e.key === "Enter" && !e.shiftKey) {
                                        e.preventDefault();
                                        handleUpdateStep(step.id, editingText);
                                    }
                                    if (e.key === "Escape")
                                        editingStepId = null;
                                }}
                            />
                            <div class="step-draft-actions">
                                <button
                                    class="step-confirm-btn"
                                    onclick={() =>
                                        handleUpdateStep(step.id, editingText)}
                                    >Save</button
                                >
                                <button
                                    class="step-cancel-btn"
                                    onclick={() => (editingStepId = null)}
                                    >Cancel</button
                                >
                            </div>
                        {:else}
                            <p class="step-text">{step.text}</p>
                            <div class="step-actions">
                                <button
                                    class="step-action-btn"
                                    onclick={() => {
                                        editingStepId = step.id;
                                        editingText = step.text;
                                    }}>Edit</button
                                >
                                <button
                                    class="step-action-btn step-action-btn--danger"
                                    onclick={() => handleDeleteStep(step.id)}
                                    >Remove</button
                                >
                            </div>
                        {/if}
                    </div>
                </li>
            {/each}

            {#if addingStep}
                <li class="step-item step-item--draft">
                    <div class="step-drag-handle" aria-hidden="true"></div>
                    <div class="step-number step-number--draft">
                        {listSteps.length + 1}
                    </div>
                    <div class="step-body">
                        <textarea
                            class="step-edit-input"
                            bind:value={draftText}
                            rows="3"
                            placeholder="Describe this step…"
                            autofocus
                            onkeydown={(e) => {
                                if (e.key === "Enter" && !e.shiftKey) {
                                    e.preventDefault();
                                    handleAddStep();
                                }
                                if (e.key === "Escape") {
                                    addingStep = false;
                                    draftText = "";
                                }
                            }}
                        />
                        <div class="step-draft-actions">
                            <button
                                class="step-confirm-btn"
                                onclick={handleAddStep}>Add</button
                            >
                            <button
                                class="step-cancel-btn"
                                onclick={() => {
                                    addingStep = false;
                                    draftText = "";
                                }}>Cancel</button
                            >
                        </div>
                    </div>
                </li>
            {/if}
        </ol>
    {/if}
</section>
