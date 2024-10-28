<script>
	import { createEventDispatcher } from "svelte";

	export let checked;
	export let id = "";
	const dispatch = createEventDispatcher();

	function handle_change(event) {
		checked = event.target.checked;
		dispatch("change", { id, checked });
	}
</script>

<label>
	<input type="checkbox" bind:checked on:change={handle_change} />
	<slot></slot>
</label>
