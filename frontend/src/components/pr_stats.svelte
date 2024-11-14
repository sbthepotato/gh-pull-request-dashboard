<script>
	import Icon from "./icon.svelte";

	export let who = "";
	export let count = 0;

	let size = "16px";
</script>

<div class="container">
	{#if who === "ready to merge"}
		<p>Ready to Merge</p>
		<Icon name="check-16" color="green" {size} />
	{:else if who === "Changes Requested"}
		<p>Changes Requested</p>
		<Icon name="file-diff-16" color="red" {size} />
	{:else if who === "total"}
		<p>Total Pull Requests</p>
		<Icon name="git-pull-request-16" color="green" {size} />
	{:else if who === "missing status"}
		<p>Missing Status</p>
		<Icon name="alert-16" color="yellow" {size} />
	{:else if who === "error"}
		<p>Error</p>
		<Icon name="alert-16" color="red" {size} />
	{:else}
		<p>Waiting on {who}</p>
		<Icon name="dot-fill-16" color="yellow" {size} />
	{/if}

	<p class="count">{count}</p>
</div>

<style>
	div.container {
		display: inline-block;
		padding: 8px;
		margin: 8px;
		border-radius: 8px;
		border: 1px solid var(--border);
	}

	p {
		margin: 4px 0;
	}

	p.count {
		font-weight: bold;
		font-size: x-large;
	}
</style>
