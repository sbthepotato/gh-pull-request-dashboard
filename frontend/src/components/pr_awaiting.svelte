<script>
	import Icon from "./icon.svelte";

	export let awaiting = "";
	export let message = "";

	let size = "16px";
</script>

<div class={awaiting} title={message}>
	{#if awaiting !== ""}
		{#if awaiting === "APPROVED"}
			<Icon name="check-16" {size} color="green" />
		{:else if awaiting === "Changes Requested"}
			<Icon name="file-diff-16" {size} color="red" />
		{:else if awaiting === "review"}
			<Icon name="dot-fill-16" {size} color="yellow" />
		{:else if awaiting === "error"}
			<Icon name="alert-16" {size} color="red" />
		{:else}
			{awaiting}
		{/if}
	{:else}
		<Icon name="alert-16" {size} color="yellow" />
	{/if}
</div>

<style>
	div {
		display: inline-block;
	}

	.error:hover {
		cursor: pointer;
	}
</style>
