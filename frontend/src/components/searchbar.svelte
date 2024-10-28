<script>
	import { createEventDispatcher } from "svelte";
	import Icon from "./icon.svelte";

	export let placeholder = "Search...";
	export let value = "";

	const dispatch = createEventDispatcher();

	function handle_change(event) {
		value = event.target.value;
		dispatch("change", { value });
	}

	function clear_search() {
		value = "";
		dispatch("change", { value });
	}
</script>

<div class="search-container">
	<input
		type="search"
		{placeholder}
		class={value ? "no-rounding" : ""}
		bind:value
		on:change={handle_change}
		on:input={handle_change} />
	{#if value}
		<button on:click={clear_search}>X</button>
	{/if}
</div>

<style>
	.search-container {
		margin: 8px auto;
		position: relative;
	}

	input {
		display: inline-block;
		background-color: var(--content-bg-alt);
		color: var(--text);
		border: none;
		box-sizing: border-box;
		padding: 8px 16px;
		border-radius: 4px;
		width: 35vw;
		min-width: 256px;
	}

	input:focus {
		outline: none;
		border: 1px solid var(--border-blue);
	}

	.no-rounding {
		border-radius: 4px 0px 0px 4px;
	}

	button {
		position: absolute;
		height: 100%;
		width: 32px;
		border: none;
		background-color: var(--border-blue);
		border-radius: 0px 4px 4px 0px;
	}

	button:hover {
		cursor: pointer;
	}

	input::-webkit-search-decoration,
	input::-webkit-search-cancel-button,
	input::-webkit-search-results-button,
	input::-webkit-search-results-decoration {
		-webkit-appearance: none;
	}
</style>
