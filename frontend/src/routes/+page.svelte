<script lang="ts">
    import '../app.css'
    import { goto } from "$app/navigation";
	import { onMount } from "svelte";

    onMount(() => goto('/login'))
</script>
