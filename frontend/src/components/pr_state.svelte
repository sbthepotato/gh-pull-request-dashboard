<script>
	import Icon from "./icon.svelte";

	export let state;

	let size = "16px";
</script>

{#if state === "open"}
	<Icon name="git-pull-request-16" color="green" {size} />
{:else if state === "draft"}
	<Icon name="git-pull-request-draft-16" color="grey" {size} />
{/if}
