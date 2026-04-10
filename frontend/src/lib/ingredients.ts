type Specifier = "g" | "kg" | "ml" | "l" | "tsp" | "tbsp" | "cup" | "oz" | "lb" | "piece";

const conversionRates: Record<Specifier, Partial<Record<Specifier, number>>> = {
	g: { kg: 0.001, oz: 0.0353, lb: 0.00220, g: 1, ml: 1, l: 0.001 },
	kg: { g: 1000, oz: 35.274, lb: 2.2046, kg: 1 },
	ml: { l: 0.001, tsp: 0.2029, tbsp: 0.0676, cup: 0.00423, ml: 1 },
	l: { ml: 1000, tsp: 202.9, tbsp: 67.628, cup: 4.2268, l: 1 },
	tsp: { tbsp: 0.3333, cup: 0.0208, ml: 4.929, tsp: 1 },
	tbsp: { tsp: 3, cup: 0.0625, ml: 14.787, tbsp: 1 },
	cup: { tsp: 48, tbsp: 16, ml: 236.588, cup: 1 },
	oz: { g: 28.3495, lb: 0.0625, oz: 1 },
	lb: { g: 453.592, oz: 16, lb: 1 },
	piece: { piece: 1 },
};

export class Ingredient {
	name: string;
	amount: number = $state(0);
	specifier: Specifier = $state("g");

	constructor(name: string, amount: number, specifier: Specifier) {
		this.name = name;
		this.amount = amount;
		this.specifier = specifier;
	}

	convertTo(target: Specifier): void {
		const rate = conversionRates[this.specifier]?.[target];
		if (rate === undefined) {
			throw new Error(`No conversion from ${this.specifier} to ${target}`);
		}
		this.amount = this.amount * rate;
		this.specifier = target;
	}
}
