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
  <nav class="flex items-center justify-between px-6 py-5 max-w-5xl mx-auto">
    <span class="font-display text-lg text-text">Owned Cloud</span>
    <a href="/login" class="text-sm text-text-muted hover:text-text transition-colors">Sign in</a>
  </nav>

  <section class="max-w-5xl mx-auto px-6 pt-10 pb-20 grid md:grid-cols-2 gap-12 items-center">
    <div>
      <p class="font-mono text-xs text-accent tracking-widest uppercase mb-4">
        Self-hosted family cloud
      </p>
      <h1 class="font-display text-4xl sm:text-5xl text-text leading-[1.1] mb-5">
        Your files. Your server. Your rules.
      </h1>
      <p class="text-text-muted text-base leading-relaxed mb-8 max-w-md">
        Owned Cloud is a private file drive you run yourself — folders, sharing, and storage
        limits, with no third party holding your family's data.
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

  <section class="max-w-5xl mx-auto px-6 pb-20">
    <div class="grid sm:grid-cols-2 gap-6">
      {#each features as feature (feature.title)}
        <div class="flex gap-3">
          <feature.icon size={18} class="text-accent-secondary shrink-0 mt-0.5" />
          <div>
            <h3 class="text-sm font-medium text-text mb-1">{feature.title}</h3>
            <p class="text-sm text-text-muted leading-relaxed">{feature.body}</p>
          </div>
        </div>
      {/each}
    </div>
  </section>

  <footer class="border-t border-border">
    <div class="max-w-5xl mx-auto px-6 py-6 flex flex-wrap items-center justify-between gap-4">
      <div class="flex gap-2">
        {#each stack as tech (tech)}
          <span
            class="font-mono text-xs text-text-muted border border-border rounded px-2 py-1"
          >
            {tech}
          </span>
        {/each}
      </div>
      <a
        href={REPO_URL}
        target="_blank"
        rel="noreferrer"
        class="text-sm text-text-muted hover:text-text transition-colors"
      >
        Source on GitHub
      </a>
    </div>
  </footer>
</div>