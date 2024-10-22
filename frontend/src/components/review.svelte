<script>
	import Icon from "./icon.svelte";

	export let state;

	let size = "16px";
</script>

{#if state === "APPROVED"}
	<Icon name="check-16" color="green" {size} />
{:else if state === "CHANGES_REQUESTED"}
	<Icon name="file-diff-16" color="red" {size} />
{:else if state === "COMMENTED"}
	<Icon name="comment-16" color="yellow" {size} />
{:else if state === "REVIEW_REQUESTED"}
	<Icon name="dot-fill-16" color="yellow" {size} />
{/if}
