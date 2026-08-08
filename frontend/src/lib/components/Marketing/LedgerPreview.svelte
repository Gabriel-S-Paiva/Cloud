<script lang="ts">
  import { Folder, FileText, Image, FileSpreadsheet } from '@lucide/svelte';

  const entries = [
    { type: 'folder', name: 'Family Photos', meta: '48 items', icon: Folder },
    { type: 'folder', name: 'Home Videos', meta: '12 items', icon: Folder },
    { type: 'file', name: 'Tax Documents 2026.pdf', meta: '2.4 MB', icon: FileText },
    { type: 'file', name: 'Reunion Panorama.jpg', meta: '18.1 MB', icon: Image },
    { type: 'file', name: 'Household Budget.xlsx', meta: '340 KB', icon: FileSpreadsheet }
  ];
</script>

<div class="rounded-lg border border-border bg-surface-raised overflow-hidden shadow-sm">
  <div class="flex items-center justify-between px-4 py-2.5 border-b border-border">
    <span class="font-mono text-xs text-text-muted tracking-wide">/home</span>
    <span class="font-mono text-xs text-text-muted">5 items</span>
  </div>
  <ul>
    {#each entries as entry, i (entry.name)}
      <li
        class="flex items-center gap-3 px-4 py-2.5 {i !== entries.length - 1
          ? 'border-b border-border'
          : ''}"
      >
        <entry.icon
          size={16}
          class={entry.type === 'folder' ? 'text-accent shrink-0' : 'text-accent-secondary shrink-0'}
        />
        <span class="text-sm text-text truncate flex-1">{entry.name}</span>
        <span class="font-mono text-xs text-text-muted shrink-0">{entry.meta}</span>
      </li>
    {/each}
  </ul>
</div>