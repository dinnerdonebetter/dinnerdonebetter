<script lang="ts">
	import { enhance } from '$app/forms';
	import {
		PageContainer,
		FormField,
		Input,
		Button,
		Alert,
		Link,
		Autocomplete,
		NumberInput,
		Card
	} from '$lib/components';
	import { createRecipeCreatorState } from '$lib/recipes/RecipeCreatorState';
	import type {
		ValidPreparation,
		ValidIngredient,
		ValidIngredientPreparation,
		ValidIngredientMeasurementUnit,
		ValidPreparationInstrument
	} from '$lib/recipes/client-types';
	import { RecipeStepProductType } from '$lib/recipes/client-enums';

	let { data, form } = $props();

	const state = $state(createRecipeCreatorState());

	function debounce<T extends (...args: unknown[]) => void>(fn: T, ms: number) {
		let timer: ReturnType<typeof setTimeout> | null = null;
		return (...args: Parameters<T>) => {
			if (timer) clearTimeout(timer);
			timer = setTimeout(() => fn(...args), ms);
		};
	}

	const debouncedPrepSearch = debounce(async (stepIndex: number, query: string) => {
		const results = await fetchPreparations(query);
		state.setPreparationSuggestions(
			stepIndex,
			results.map((r) => ({ id: r.id, name: r.label }))
		);
	}, 300);

	const debouncedIngredientSearch = debounce(
		async (stepIndex: number, ingIdx: number, query: string) => {
			const prep = state.stepHelpers[stepIndex]?.selectedPreparation;
			if (!prep) return;
			const results = await fetchIngredients(query, prep.id);
			state.setIngredientSuggestions(stepIndex, ingIdx, results);
		},
		300
	);

	async function fetchPreparations(query: string): Promise<{ id: string; label: string }[]> {
		if (query.length < 2) return [];
		const res = await fetch(`/api/recipes/search-preparations?q=${encodeURIComponent(query)}`);
		const json = await res.json();
		return (json.results ?? []).map((p: ValidPreparation) => ({ id: p.id, label: p.name }));
	}

	async function fetchIngredients(
		query: string,
		preparationId: string
	): Promise<ValidIngredient[]> {
		if (query.length < 2 || !preparationId) return [];
		const res = await fetch(
			`/api/recipes/search-ingredients?q=${encodeURIComponent(query)}&preparationId=${encodeURIComponent(preparationId)}`
		);
		const json = await res.json();
		return json.results ?? [];
	}

	async function fetchIngredientPreparations(
		preparationId: string
	): Promise<ValidIngredientPreparation[]> {
		const res = await fetch(
			`/api/recipes/ingredient-preparations?preparationId=${encodeURIComponent(preparationId)}`
		);
		const json = await res.json();
		return json.results ?? [];
	}

	const productTypeOptions = [
		{ value: '0', label: 'Ingredient' },
		{ value: '1', label: 'Instrument' },
		{ value: '2', label: 'Vessel' }
	];

	function slugify(s: string): string {
		return s
			.toLowerCase()
			.trim()
			.replace(/[^\w\s-]/g, '')
			.replace(/[\s_-]+/g, '-')
			.replace(/^-+|-+$/g, '');
	}

	$effect(() => {
		const name = state.recipe.name;
		if (name && !state.recipe.slug) {
			state.updateRecipeField('slug', slugify(name));
		}
	});
</script>

<PageContainer>
	<h1>Create Recipe</h1>
	<p><Link href="/meal_plans">Back to Meal Plans</Link></p>

	{#if form?.error}
		<Alert variant="error">{form.error}</Alert>
	{/if}
	{#if state.submissionError}
		<Alert variant="error">{state.submissionError}</Alert>
	{/if}

	<form
		method="POST"
		action="?/default"
		use:enhance
		class="recipe-form"
	>
		<input type="hidden" name="recipe" value={JSON.stringify(state.recipe)} />

		<section class="recipe-section">
			<h2>Recipe Details</h2>
			<FormField id="name" label="Name" required>
				<Input
					id="name"
					bind:value={state.recipe.name}
					placeholder="Recipe name"
					required
				/>
			</FormField>
			<FormField id="slug" label="URL Slug">
				<Input
					id="slug"
					bind:value={state.recipe.slug}
					placeholder="recipe-url-slug"
				/>
			</FormField>
			<FormField id="source" label="Source">
				<Input
					id="source"
					bind:value={state.recipe.source}
					placeholder="Recipe source"
				/>
			</FormField>
			<FormField id="description" label="Description">
				<textarea
					id="description"
					name="description"
					class="textarea"
					placeholder="Recipe description"
					bind:value={state.recipe.description}
					rows="4"
				></textarea>
			</FormField>
			<FormField id="portions" label="Estimated Portions" required>
				<NumberInput
					id="portions"
					bind:value={state.recipe.estimatedPortions!.min}
					min={1}
					required
				/>
			</FormField>
			<FormField id="portionName" label="Portion Name">
				<Input
					id="portionName"
					bind:value={state.recipe.portionName}
					placeholder="portion"
				/>
			</FormField>
			<FormField id="pluralPortionName" label="Plural Portion Name">
				<Input
					id="pluralPortionName"
					bind:value={state.recipe.pluralPortionName}
					placeholder="portions"
				/>
			</FormField>
		</section>

		<section class="recipe-section">
			<h2>Steps</h2>
			{#each state.recipe.steps as step, stepIndex}
				<Card
					title="Step {stepIndex + 1}"
					collapsible
					bind:expanded={state.stepHelpers[stepIndex].show}
				>
					<FormField label="Preparation" required>
							<Autocomplete
								bind:value={state.stepHelpers[stepIndex].preparationQuery}
							placeholder="Search preparation (e.g. dice, chop)"
							suggestions={state.stepHelpers[stepIndex].preparationSuggestions.map((p) => ({
								id: p.id,
								label: p.name
							}))}
							onInput={(query: string) => debouncedPrepSearch(stepIndex, query)}
							onSelect={async (item) => {
								const prep = state.stepHelpers[stepIndex].preparationSuggestions.find(
									(p) => p.id === item.id
								) ?? { id: item.id, name: item.label };
								state.setPreparation(stepIndex, prep as ValidPreparation);
								const res = await fetch(
									`/api/recipes/search-instruments?preparationId=${encodeURIComponent(prep.id)}`
								);
								const json = await res.json();
								state.setInstrumentSuggestions(stepIndex, json.results ?? []);
							}}
						/>
					</FormField>

					{#if state.stepHelpers[stepIndex].selectedPreparation}
						<div class="step-section">
							<h3>Instruments</h3>
							{#each step.instruments as instrument, instIdx}
								<div class="step-item">
									<FormField label="Instrument {instIdx + 1}">
										<Autocomplete
											value={instrument.name}
											placeholder="Search instrument"
											suggestions={state.stepHelpers[stepIndex].instrumentSuggestions.map(
												(vpi) => ({ id: vpi.id, label: vpi.instrument?.name ?? vpi.id })
											)}
											onSelect={(item) => {
												const vpi = state.stepHelpers[stepIndex].instrumentSuggestions.find(
													(x) => x.id === item.id
												);
												if (vpi) state.setInstrument(stepIndex, instIdx, vpi);
											}}
										/>
									</FormField>
									<Button
										type="button"
										variant="default"
										onclick={() => state.removeInstrumentFromStep(stepIndex, instIdx)}
									>
										Remove
									</Button>
								</div>
							{/each}
							<Button
								type="button"
								variant="default"
								onclick={() => state.addInstrumentToStep(stepIndex)}
							>
								Add Instrument
							</Button>
						</div>

						<div class="step-section">
							<h3>Ingredients</h3>
							{#each step.ingredients as ingredient, ingIdx}
								<div class="step-item">
									<FormField label="Ingredient {ingIdx + 1}">
										<Autocomplete
											bind:value={state.stepHelpers[stepIndex].ingredientQueries[ingIdx]}
											placeholder="Search ingredient"
											suggestions={(state.stepHelpers[stepIndex].ingredientSuggestions[ingIdx] ?? []).map(
												(ing) => ({ id: ing.id, label: ing.name })
											)}
											onInput={(query: string) =>
												debouncedIngredientSearch(stepIndex, ingIdx, query)}
											onSelect={async (item) => {
												const vips = await fetchIngredientPreparations(
													state.stepHelpers[stepIndex].selectedPreparation!.id
												);
												const vip = vips.find((v) => v.ingredient?.id === item.id);
												const vimus = await fetch(
													`/api/recipes/search-measurement-units?ingredientId=${encodeURIComponent(item.id)}`
												).then((r) => r.json());
												const vimuList: ValidIngredientMeasurementUnit[] = vimus.results ?? [];
												state.setIngredientMeasurementUnitSuggestions(stepIndex, ingIdx, vimuList);
												const ing = {
													id: item.id,
													name: item.label,
													validIngredientPreparationId: vip?.id
												};
												const firstVimu = vimuList[0] ?? null;
												state.setIngredient(
													stepIndex,
													ingIdx,
													ing,
													firstVimu
												);
											}}
										/>
									</FormField>
									{#if (state.stepHelpers[stepIndex].ingredientMeasurementUnitSuggestions[ingIdx] ?? []).length > 0}
										<FormField label="Measurement Unit" required>
											<select
												class="select-input"
												value={ingredient.validIngredientMeasurementUnitId ?? ''}
												onchange={(e) => {
													const id = (e.currentTarget as HTMLSelectElement).value;
													const vimu = (state.stepHelpers[stepIndex].ingredientMeasurementUnitSuggestions[ingIdx] ?? []).find(
														(v) => v.id === id
													);
													if (vimu) state.setIngredientMeasurementUnit(stepIndex, ingIdx, vimu);
												}}
											>
												<option value="" disabled>Select unit</option>
												{#each (state.stepHelpers[stepIndex].ingredientMeasurementUnitSuggestions[ingIdx] ?? []) as vimu}
													<option value={vimu.id}>
														{vimu.measurementUnit?.name ?? vimu.id}
													</option>
												{/each}
											</select>
										</FormField>
									{/if}
									<FormField label="Quantity">
										<NumberInput
											bind:value={ingredient.quantity!.min}
											min={0}
											step={0.25}
										/>
									</FormField>
									<Button
										type="button"
										variant="default"
										onclick={() => state.removeIngredientFromStep(stepIndex, ingIdx)}
									>
										Remove
									</Button>
								</div>
							{/each}
							<Button
								type="button"
								variant="default"
								onclick={() => state.addIngredientToStep(stepIndex)}
							>
								Add Ingredient
							</Button>
						</div>

						<div class="step-section">
							<h3>Products</h3>
							{#each step.products as product, prodIdx}
								<div class="step-item">
									<FormField label="Product {prodIdx + 1}">
										<Input
											bind:value={product.name}
											placeholder="Product name"
										/>
									</FormField>
									<FormField label="Type">
										<select
											class="select-input"
											value={String(product.type)}
											onchange={(e) => {
												const v = (e.currentTarget as HTMLSelectElement).value;
												state.updateProduct(stepIndex, prodIdx, {
													type: Number(v) as RecipeStepProductType
												});
											}}
										>
											{#each productTypeOptions as opt}
												<option value={opt.value}>{opt.label}</option>
											{/each}
										</select>
									</FormField>
									<Button
										type="button"
										variant="default"
										onclick={() => state.removeProductFromStep(stepIndex, prodIdx)}
									>
										Remove
									</Button>
								</div>
							{/each}
							<Button
								type="button"
								variant="default"
								onclick={() => state.addProductToStep(stepIndex)}
							>
								Add Product
							</Button>
						</div>

						<FormField label="Instructions">
							<textarea
								class="textarea"
								placeholder="Step instructions"
								bind:value={step.explicitInstructions}
								rows="3"
							></textarea>
						</FormField>
					{/if}

					<div class="step-actions">
						<Button
							type="button"
							variant="default"
							onclick={() => state.removeStep(stepIndex)}
							disabled={state.recipe.steps.length <= 1}
						>
							Remove Step
						</Button>
					</div>
				</Card>
			{/each}

			<Button type="button" variant="default" onclick={() => state.addStep()}>
				Add Step
			</Button>
		</section>

		<div class="form-actions">
			<Button type="submit">Create Recipe</Button>
		</div>
	</form>
</PageContainer>

<style>
	.recipe-form {
		display: flex;
		flex-direction: column;
		gap: var(--space-lg);
		max-width: var(--content-max-width);
	}

	.recipe-section {
		padding: var(--space-md);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		background: var(--color-surface);
	}

	.recipe-section h2 {
		margin: 0 0 var(--space-md);
		font-size: 1.25rem;
	}

	.step-section {
		margin-top: var(--space-md);
	}

	.step-section h3 {
		margin: 0 0 var(--space-sm);
		font-size: 1rem;
	}

	.step-item {
		display: flex;
		flex-wrap: wrap;
		gap: var(--space-md);
		align-items: flex-end;
		margin-bottom: var(--space-md);
	}

	.step-item :global(.form-field) {
		flex: 1;
		min-width: 12rem;
	}

	.step-actions {
		margin-top: var(--space-md);
	}

	.form-actions {
		margin-top: var(--space-md);
	}

	.textarea {
		width: 100%;
		padding: var(--space-sm);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-sm);
		font-family: var(--font-sans);
		font-size: 1rem;
		color: var(--color-text);
		resize: vertical;
	}

	.textarea:focus {
		outline: none;
		border-color: var(--color-primary);
	}

	.select-input {
		width: 100%;
		padding: var(--space-sm);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-sm);
		font-family: var(--font-sans);
		font-size: var(--font-size-base);
		color: var(--color-text);
		background: var(--color-surface);
		cursor: pointer;
		box-shadow: var(--shadow-sm);
		transition: border-color var(--transition-fast), box-shadow var(--transition-fast);
	}

	.select-input:focus {
		outline: none;
		border-color: var(--color-primary);
		box-shadow: 0 0 0 2px var(--color-primary-muted);
	}
</style>
