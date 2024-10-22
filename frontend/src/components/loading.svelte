<script>
	import Icon from "./icon.svelte";

	export let text = "loading...";
	export let size = "128px";
</script>

<div>
	<p>{text}</p>
	<Icon name="mark-github-24" color="rainbow" {size} />
</div>

<style>
</style>
