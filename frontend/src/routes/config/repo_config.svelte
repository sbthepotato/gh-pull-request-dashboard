<script>
	import { onMount } from "svelte";
	import Button from "../../components/button.svelte";

	let repos = [];
	let result = "";
	let err = "";

	async function get_repos(refresh) {
		try {
			repos = [];
			err = "";

			let url = "api/config/get_repos";

			if (refresh) {
				url = url + "?refresh=y";
			}

			const response = await fetch(url);

			if (response.ok) {
				repos = await response.json();
			} else {
				throw new Error(await response.text());
			}
		} catch (error) {
			err = error.message;
		}
	}
</script>

<h2>Repository Configuration</h2>
<h3>Set the active repositories</h3>

<Button color="green" on_click={() => get_repos(true)}>
	hard refresh repository list
</Button>

{#if err !== ""}
	<p>
		{err}
	</p>
{:else if repos.length > 0}
	<p>{repos.length} repos found</p>
	<table>
		<tbody>
			{#each repos as repo}
				<tr>
					<td>
						{repo.name}
					</td>
				</tr>
			{/each}
		</tbody>
	</table>
{/if}
