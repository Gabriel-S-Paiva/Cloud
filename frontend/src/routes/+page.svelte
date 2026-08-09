<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { auth } from '$lib/stores/auth.svelte';
	import Button from '$lib/components/UI/Button/Button.svelte';
	import LedgerPreview from '$lib/components/Marketing/LedgerPreview.svelte';
	import { Share2, Gauge, ShieldCheck, Server } from '@lucide/svelte';

	const REPO_URL = 'https://github.com/Gabriel-S-Paiva/Cloud';

	const features = [
		{
			icon: Share2,
			title: 'Folders & sharing',
			body: 'Organize files into folders and share individual items with view or edit access.'
		},
		{
			icon: Gauge,
			title: 'Storage quotas',
			body: "Every account gets a set amount of space, so one user can't fill the whole drive."
		},
		{
			icon: ShieldCheck,
			title: 'Admin-approved accounts',
			body: 'New sign-ups wait for approval before they can access anything.'
		},
		{
			icon: Server,
			title: 'Runs on your hardware',
			body: 'No subscription, no third-party server. Docker, behind Caddy, on a machine you control.'
		}
	];

	const stack = ['Go', 'SvelteKit', 'SQLite', 'Caddy'];
</script>

<div class="min-h-screen bg-surface">
	<nav class="mx-auto flex max-w-5xl items-center justify-between px-6 py-5">
		<span class="font-display text-lg text-text">Owned Cloud</span>
		<a href="/login" class="text-sm text-text-muted transition-colors hover:text-text">Sign in</a>
	</nav>

	<section class="mx-auto grid max-w-5xl items-center gap-12 px-6 pt-10 pb-20 md:grid-cols-2">
		<div>
			<p class="mb-4 font-mono text-xs tracking-widest text-accent uppercase">
				Self-hosted family cloud
			</p>
			<h1 class="mb-5 font-display text-4xl leading-[1.1] text-text sm:text-5xl">
				Your files. Your server. Your rules.
			</h1>
			<p class="mb-8 max-w-md text-base leading-relaxed text-text-muted">
				Owned Cloud is a private file drive you run yourself — folders, sharing, and storage limits,
				with no third party holding your family's data.
			</p>
			<div class="flex items-center gap-3">
				<Button onclick={() => goto('/login')}>Sign in</Button>
				<a href={REPO_URL} target="_blank" rel="noreferrer">
					<Button variant="secondary">View source</Button>
				</a>
			</div>
		</div>

		<LedgerPreview />
	</section>

	<section class="mx-auto max-w-5xl px-6 pb-20">
		<div class="grid gap-6 sm:grid-cols-2">
			{#each features as feature (feature.title)}
				<div class="flex gap-3">
					<feature.icon size={18} class="mt-0.5 shrink-0 text-accent-secondary" />
					<div>
						<h3 class="mb-1 text-sm font-medium text-text">{feature.title}</h3>
						<p class="text-sm leading-relaxed text-text-muted">{feature.body}</p>
					</div>
				</div>
			{/each}
		</div>
	</section>

	<footer class="border-t border-border">
		<div class="mx-auto flex max-w-5xl flex-wrap items-center justify-between gap-4 px-6 py-6">
			<div class="flex gap-2">
				{#each stack as tech (tech)}
					<span class="rounded border border-border px-2 py-1 font-mono text-xs text-text-muted">
						{tech}
					</span>
				{/each}
			</div>
			<a
				href={REPO_URL}
				target="_blank"
				rel="noreferrer"
				class="text-sm text-text-muted transition-colors hover:text-text"
			>
				Source on GitHub
			</a>
		</div>
	</footer>
</div>
