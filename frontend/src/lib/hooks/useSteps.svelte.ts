import { PUBLIC_SERVER_URL, PUBLIC_SERVER_PORT } from "$env/static/public";

const SERVER_URL = PUBLIC_SERVER_PORT
	? `${PUBLIC_SERVER_URL}:${PUBLIC_SERVER_PORT}`
	: PUBLIC_SERVER_URL;

export function useSteps(id: number) {
	let listSteps: Array<{ id: number; text: string; stepNumber: number }> = $state([]);

	$effect(() => {
		const get = async () => {
			const response = await fetch(`${SERVER_URL}/steps/${id}`);
			if (!response.ok) return;
			const data = await response.json();
			if (data.steps != null) {
				listSteps = data.steps.map((s: { step_id: number; step_text: string; step_number: number }) => ({
					id: s.step_id,
					text: s.step_text,
					stepNumber: s.step_number,
				}));
			}
		};
		get();
	});

	return {
		get listSteps() { return listSteps; },
		set listSteps(v) { listSteps = v; },
	};
}
