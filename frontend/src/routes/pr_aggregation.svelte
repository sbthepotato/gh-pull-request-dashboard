<script>
	import PRStats from "../components/pr_stats.svelte";

	export let pr_stats;
</script>

{#if pr_stats != undefined}
	{#each Object.entries(pr_stats) as [who, count]}
		<PRStats {who} {count} />
	{/each}
{/if}
