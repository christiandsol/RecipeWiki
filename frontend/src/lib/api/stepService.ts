import { PUBLIC_SERVER_URL, PUBLIC_SERVER_PORT } from "$env/static/public";
import type { GetStepResult } from "$lib/types";

const SERVER_URL = PUBLIC_SERVER_PORT
	? `${PUBLIC_SERVER_URL}:${PUBLIC_SERVER_PORT}`
	: PUBLIC_SERVER_URL;

export const addStep = async (
	recipeId: number,
	text: string
): Promise<GetStepResult | null> => {
	if (!text.trim()) return;
	const response = await fetch(`${SERVER_URL}/step`, {
		method: "POST",
		headers: { "Content-Type": "application/json" },
		body: JSON.stringify({
			recipe_id: recipeId,
			step_text: text.trim(),
		}),
	});
	if (!response.ok) return null;
	return response.json();
};

export const updateStep = async (step_id: number, newText: string): Promise<boolean> => {
	const trimmedText = newText.trim()
	if (!trimmedText) return false;
	const response = await fetch(`${SERVER_URL}/step`, {
		method: "PATCH",
		headers: { "Content-Type": "application/json" },
		body: JSON.stringify({ id: step_id, text: trimmedText }),
	});
	return response.ok
};

export const deleteStep = async (stepId: number) => {
	const response = await fetch(`${SERVER_URL}/step`, {
		method: "DELETE",
		headers: { "Content-Type": "application/json" },
		body: JSON.stringify({ id: stepId }),
	});
	if (response.ok) {
		listSteps = listSteps.filter((s) => s.id !== stepId);
	}
};

export const reorderStep = async (
	stepId: number,
	before: number,
	after: number,
) => {
	const response = await fetch(`${SERVER_URL}/step/reorder`, {
		method: "PATCH",
		headers: { "Content-Type": "application/json" },
		body: JSON.stringify({ step_id: stepId, before, after }),
	});
	if (!response.ok) {
		console.error("[ERROR] Failed to reorder step");
	}
};

