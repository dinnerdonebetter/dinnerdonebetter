<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		narrow?: boolean;
		class?: string;
		children?: Snippet;
	}

	let { narrow = false, class: className = '', children }: Props = $props();
</script>

<main class="page-container {narrow ? 'page-container--narrow' : ''} {className}">
	{#if children}
		{@render children()}
	{/if}
</main>

<style>
	.page-container {
		max-width: var(--content-max-width);
		margin: var(--space-lg) auto;
		padding: var(--space-md);
	}

	.page-container--narrow {
		max-width: var(--content-max-width-narrow);
	}
</style>
