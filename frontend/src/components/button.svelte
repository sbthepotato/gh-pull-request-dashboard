<script>
	import { goto } from "$app/navigation";

	export let type = "button";
	export let onClick = () => {};
	export let to = null;

	const url_prefix = import.meta.env.VITE_URL_PATH;

	function handleClick(event) {
		if (to) {
			goto(url_prefix + to);
		} else {
			onClick(event);
		}
	}
</script>

<button {type} on:click={handleClick}>
	<slot></slot>
</button>

<style>
	button {
		display: inline;
	}
</style>
