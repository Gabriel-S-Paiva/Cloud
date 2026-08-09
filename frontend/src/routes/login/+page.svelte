<script lang="ts">
	import { auth } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';
	import Input from '$lib/components/UI/Input/Input.svelte';
	import { User, Lock, Eye, EyeOff } from '@lucide/svelte';
	import { toast } from '$lib/stores/toast.svelte';
	import Toast from '$lib/components/UI/Toast/Toast.svelte';

	let mode = $state<'login' | 'register'>('login');

	let loginUsername = $state('');
	let loginPassword = $state('');
	let PasswordVisible = $state(false);
	let loginError = $state<string | null>(null);

	let registerUsername = $state('');
	let registerPassword = $state('');

	async function handleLogin(e: Event) {
		e.preventDefault();
		loginError = null;
		try {
			await auth.login(loginUsername, loginPassword);
			toast.success('Login successfull.');
			goto('/home');
		} catch (err) {
			err instanceof Error ? toast.error(err.message) : toast.error('Login failed');
		}
	}

	async function handleRegister(e: Event) {
		e.preventDefault();
		try {
			await auth.register(registerUsername, registerPassword);
			toast.success('Account created - waiting for admin approval.');
		} catch (err) {
			err instanceof Error ? toast.error(err.message) : toast.error('Registration failed');
		}
	}
</script>

<div class="flex min-h-screen items-center justify-center bg-surface px-4">
	<div class="w-full max-w-sm">
		<h1 class="mb-6 text-center font-display text-2xl text-text">Owned Cloud</h1>

		<div class="mb-6 flex border-b border-border">
			<button
				class="flex-1 pb-3 text-sm font-medium transition-colors {mode === 'login'
					? 'border-b-2 border-accent text-accent'
					: 'text-text-muted'}"
				onclick={() => (mode = 'login')}
			>
				Login
			</button>
			<button
				class="flex-1 pb-3 text-sm font-medium transition-colors {mode === 'register'
					? 'border-b-2 border-accent text-accent'
					: 'text-text-muted'}"
				onclick={() => (mode = 'register')}
			>
				Register
			</button>
		</div>

		{#if mode === 'login'}
			<form onsubmit={handleLogin} class="flex flex-col gap-3">
				<Input bind:value={loginUsername} placeholder="username">
					{#snippet left()}<User size={16} />{/snippet}
				</Input>
				<Input
					bind:value={loginPassword}
					placeholder="password"
					type={PasswordVisible ? 'text' : 'password'}
				>
					{#snippet left()}<Lock size={16} />{/snippet}
					{#snippet right()}
						<button type="button" onclick={() => (PasswordVisible = !PasswordVisible)}>
							{#if PasswordVisible}<EyeOff size={16} />{:else}<Eye size={16} />{/if}
						</button>
					{/snippet}
				</Input>
				<button
					type="submit"
					class="mt-1 h-10 rounded-md bg-accent text-sm font-medium text-surface"
				>
					Login
				</button>
				{#if loginError}<p class="text-sm text-danger">{loginError}</p>{/if}
			</form>
		{:else}
			<form onsubmit={handleRegister} class="flex flex-col gap-3">
				<Input bind:value={registerUsername} placeholder="username">
					{#snippet left()}<User size={16} />{/snippet}
				</Input>
				<Input
					bind:value={registerPassword}
					placeholder="password"
					type={PasswordVisible ? 'text' : 'password'}
				>
					{#snippet left()}<Lock size={16} />{/snippet}
					{#snippet right()}
						<button type="button" onclick={() => (PasswordVisible = !PasswordVisible)}>
							{#if PasswordVisible}<EyeOff size={16} />{:else}<Eye size={16} />{/if}
						</button>
					{/snippet}
				</Input>
				<button
					type="submit"
					class="mt-1 h-10 rounded-md bg-accent text-sm font-medium text-surface"
				>
					Register
				</button>
			</form>
		{/if}
	</div>
	<Toast />
</div>
