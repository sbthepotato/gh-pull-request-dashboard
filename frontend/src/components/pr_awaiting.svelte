<script>
	import Icon from "./icon.svelte";

	export let awaiting = "";
</script>

{#if awaiting != undefined}
	{#if awaiting === "APPROVED"}
		<Icon name="check-16" height="16px" width="16px" color="green" />
	{:else if awaiting === "Changes Requested"}
		<Icon name="file-diff-16" height="16px" width="16px" color="red" />
	{:else if awaiting === "Review"}
		<Icon name="dot-fill-16" height="16px" width="16px" color="yellow" />
	{:else}
		{awaiting}
	{/if}
{:else}
	<Icon name="alert-16" height="16px" width="16px" color="red" />
{/if}
