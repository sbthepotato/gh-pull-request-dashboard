<script>
	import { goto } from "$app/navigation";

	export let type = "button";
	export let onClick = () => {};
	export let to = null;
	export let color = "grey";

	const url_prefix = import.meta.env.VITE_URL_PATH;

	function handleClick(event) {
		if (to) {
			goto(url_prefix + to);
		} else {
			onClick(event);
		}
	}
</script>

<button class={color} {type} on:click={handleClick}>
	<slot></slot>
</button>

<style>
	button {
		color: var(--text);
		display: inline;
		border: none;
		padding: 8px;
		border-radius: 8px;
		font-weight: bold;
		margin: 4px;
	}

	button.grey {
		background-color: var(--button-grey);
	}

	button.green {
		background-color: var(--button-green);
	}

	button:hover {
		cursor: pointer;
	}
</style>
