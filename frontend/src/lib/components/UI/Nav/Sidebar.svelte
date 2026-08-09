<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { auth } from '$lib/stores/auth.svelte';
	import {
		Cloud,
		FolderOpen,
		Share2,
		CircleUser,
		ShieldCheck,
		ChevronsLeft,
		ChevronsRight,
		LogOut,
		Plus
	} from '@lucide/svelte';

	let collapsed = $state(false);

	const links = $derived([
		{ href: '/home/' as const, label: 'Files', icon: FolderOpen },
		{ href: '/shares/incoming' as const, label: 'Shared With Me', icon: Share2 },
		{ href: '/shares/outgoing' as const, label: 'Manage Shares', icon: Share2 },
		{ href: '/profile' as const, label: 'Profile', icon: CircleUser },
		...(auth.user?.role === 'Admin' ? [{ href: '/admin' as const, label: 'Admin', icon: ShieldCheck }] : [])
	]);

	const isActive = (href: string) => page.url.pathname === href;

	const logout = async () => {
		await auth.logout();
		goto(resolve('/login'));
	};
</script>

<aside
	class="flex h-screen flex-col border-r border-border bg-surface-raised transition-[width] duration-150 {collapsed
		? 'w-16'
		: 'w-56'}"
>
	<div class="flex h-16 shrink-0 items-center gap-2 px-4">
		<Cloud size={22} class="shrink-0 text-accent" />
		{#if !collapsed}
			<span class="truncate font-display text-lg text-text">Owned Cloud</span>
		{/if}
	</div>

	<div class="mb-4 px-3">
		<button
			onclick={() => goto(resolve('/home/'))}
			class="flex h-10 w-full items-center justify-center gap-2 rounded-md bg-accent text-sm font-medium text-surface"
		>
			<Plus size={16} />
			{#if !collapsed}New{/if}
		</button>
	</div>

	<nav class="flex flex-1 flex-col gap-1 px-2">
		{#each links as { href, label, icon: Icon } (href)}
			<button
				onclick={() => goto(resolve(href))}
				class="relative flex h-10 items-center gap-3 rounded-md px-3 text-sm transition-colors {isActive(
					href
				)
					? 'bg-accent/10 text-accent'
					: 'text-text-muted hover:bg-surface hover:text-text'}"
			>
				{#if isActive(href)}
					<span class="absolute top-1.5 bottom-1.5 left-0 w-0.5 rounded-full bg-accent"></span>
				{/if}
				<Icon size={18} class="shrink-0" />
				{#if !collapsed}<span class="truncate">{label}</span>{/if}
			</button>
		{/each}
	</nav>

	<div class="flex flex-col gap-1 border-t border-border px-2 pt-3 pb-3">
		<button
			onclick={logout}
			class="flex h-10 items-center gap-3 rounded-md px-3 text-sm text-text-muted transition-colors hover:text-danger"
		>
			<LogOut size={18} class="shrink-0" />
			{#if !collapsed}Logout{/if}
		</button>
		<button
			onclick={() => (collapsed = !collapsed)}
			class="flex h-10 items-center gap-3 rounded-md px-3 text-sm text-text-muted transition-colors hover:text-text"
		>
			{#if collapsed}
				<ChevronsRight size={18} />
			{:else}
				<ChevronsLeft size={18} />
				<span>Collapse</span>
			{/if}
		</button>
	</div>
</aside>
