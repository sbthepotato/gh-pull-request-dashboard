<script>
	import { createEventDispatcher } from "svelte";

	export let placeholder = "Search...";
	export let value = "";

	const dispatch = createEventDispatcher();

	function handle_change(event) {
		value = event.target.value;
		dispatch("change", { value });
	}
</script>

<input type="text" {placeholder} bind:value on:change={handle_change} />

<style>
	input {
		display: block;
		margin: 8px auto;
		background-color: var(--content-bg-alt);
		color: var(--text);
		border: 1px solid var(--blue);
		box-sizing: border-box;
		padding: 8px 16px;
		border-radius: 4px;
		width: 35vw;
		min-width: 256px;
	}
</style>
