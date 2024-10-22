<script>
	import Icon from "./icon.svelte";

	export let awaiting = "";

	let size = "16px";
</script>

{#if awaiting !== ""}
	{#if awaiting === "APPROVED"}
		<Icon name="check-16" {size} color="green" />
	{:else if awaiting === "Changes Requested"}
		<Icon name="file-diff-16" {size} color="red" />
	{:else if awaiting === "review"}
		<Icon name="dot-fill-16" {size} color="yellow" />
	{:else}
		{awaiting}
	{/if}
{:else}
	<Icon name="alert-16" {size} color="red" />
{/if}
