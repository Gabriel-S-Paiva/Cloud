<script lang="ts">
	import { goto } from '$app/navigation';
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
		{ href: '/home', label: 'Files', icon: FolderOpen },
		{ href: '/shares/incoming', label: 'Shared With Me', icon: Share2 },
		{ href: '/shares/outgoing', label: 'Manage Shares', icon: Share2 },
		{ href: '/profile', label: 'Profile', icon: CircleUser },
		...(auth.user?.role === 'Admin' ? [{ href: '/admin', label: 'Admin', icon: ShieldCheck }] : [])
	]);

	const isActive = (href: string) => page.url.pathname === href;

	const logout = async () => {
		await auth.logout();
		goto('/login');
	};
</script>

<aside
	class="h-screen border-r border-border bg-surface-raised flex flex-col transition-[width] duration-150 {collapsed ? 'w-16' : 'w-56'}"
>
	<div class="flex items-center gap-2 px-4 h-16 shrink-0">
		<Cloud size={22} class="text-accent shrink-0" />
		{#if !collapsed}
			<span class="font-display text-lg text-text truncate">Owned Cloud</span>
		{/if}
	</div>

	<div class="px-3 mb-4">
		<button
			onclick={() => goto('/home')}
			class="w-full h-10 rounded-md bg-accent text-surface flex items-center justify-center gap-2 text-sm font-medium"
		>
			<Plus size={16} />
			{#if !collapsed}New{/if}
		</button>
	</div>

	<nav class="flex-1 px-2 flex flex-col gap-1">
		{#each links as { href, label, icon: Icon } (href)}
			<button
				onclick={() => goto(href)}
				class="relative h-10 rounded-md flex items-center gap-3 px-3 text-sm transition-colors {isActive(href) ? 'text-accent bg-accent/10' : 'text-text-muted hover:text-text hover:bg-surface'}"
			>
				{#if isActive(href)}
					<span class="absolute left-0 top-1.5 bottom-1.5 w-0.5 bg-accent rounded-full"></span>
				{/if}
				<Icon size={18} class="shrink-0" />
				{#if !collapsed}<span class="truncate">{label}</span>{/if}
			</button>
		{/each}
	</nav>

	<div class="px-2 pb-3 flex flex-col gap-1 border-t border-border pt-3">
		<button
			onclick={logout}
			class="h-10 rounded-md flex items-center gap-3 px-3 text-sm text-text-muted hover:text-danger transition-colors"
		>
			<LogOut size={18} class="shrink-0" />
			{#if !collapsed}Logout{/if}
		</button>
		<button
			onclick={() => (collapsed = !collapsed)}
			class="h-10 rounded-md flex items-center gap-3 px-3 text-sm text-text-muted hover:text-text transition-colors"
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