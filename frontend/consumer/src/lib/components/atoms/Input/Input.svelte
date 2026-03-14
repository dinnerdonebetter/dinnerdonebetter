<script lang="ts">
	interface Props {
		id?: string;
		name?: string;
		type?: string;
		value?: string;
		placeholder?: string;
		required?: boolean;
		autocomplete?: string;
		disabled?: boolean;
		class?: string;
	}

	let {
		id,
		name,
		type = 'text',
		value = $bindable(''),
		placeholder,
		required = false,
		autocomplete,
		disabled = false,
		class: className = ''
	}: Props = $props();
</script>

<input
	{id}
	{name}
	{type}
	{placeholder}
	{required}
	autocomplete={autocomplete as HTMLInputElement['autocomplete']}
	{disabled}
	class="input {className}"
	{value}
/>

<style>
	.input {
		width: 100%;
		padding: var(--space-sm);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-sm);
		font-family: var(--font-sans);
		font-size: var(--font-size-base);
		color: var(--color-text);
		box-shadow: var(--shadow-sm);
		transition: border-color var(--transition-fast), box-shadow var(--transition-fast);
	}

	.input::placeholder {
		color: var(--color-text-muted);
	}

	.input:focus {
		outline: none;
		border-color: var(--color-primary);
		box-shadow: 0 0 0 2px var(--color-primary-muted);
	}
</style>
