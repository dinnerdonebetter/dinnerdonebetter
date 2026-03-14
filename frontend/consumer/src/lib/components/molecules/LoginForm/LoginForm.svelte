<script lang="ts">
	import type { Snippet } from 'svelte';
	import { enhance } from '$app/forms';
	import FormField from '../FormField/FormField.svelte';
	import Input from '../../atoms/Input/Input.svelte';
	import Button from '../../atoms/Button/Button.svelte';
	import Alert from '../../layout/Alert/Alert.svelte';

	interface Props {
		action?: string;
		method?: 'POST' | 'GET';
		username?: string;
		showTotp?: boolean;
		error?: string;
		submitLabel?: string;
		passkeySlot?: Snippet;
	}

	let {
		action = '?/login',
		method = 'POST',
		username = '',
		showTotp = false,
		error,
		submitLabel = 'Sign In',
		passkeySlot
	}: Props = $props();
</script>

<form
	{method}
	{action}
	use:enhance={() => {
		return async ({ result }) => {
			if (result.type === 'redirect') {
				window.location.href = result.location;
			}
		};
	}}
>
	{#if error}
		<Alert variant="error">{error}</Alert>
	{/if}
	<FormField id="username" label="Username" required>
		<Input
			id="username"
			name="username"
			type="text"
			autocomplete="username"
			value={username}
			required
			dataTestId="login-username"
		/>
	</FormField>
	<FormField id="password" label="Password" required>
		<Input
			id="password"
			name="password"
			type="password"
			autocomplete="current-password"
			required
			dataTestId="login-password"
		/>
	</FormField>
	{#if showTotp}
		<FormField id="totpToken" label="TOTP (if enabled)">
			<Input
				id="totpToken"
				name="totpToken"
				type="text"
				autocomplete="one-time-code"
				placeholder="Optional"
				dataTestId="login-totp-token"
			/>
		</FormField>
	{/if}
	<Button type="submit">{submitLabel}</Button>
	{#if passkeySlot}
		{@render passkeySlot()}
	{/if}
</form>
