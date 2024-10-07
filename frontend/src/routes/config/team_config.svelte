<script>
	import { onMount } from "svelte";
	import Button from "../../components/button.svelte";

	let teams = [];
	let result = "";
	let err = "";

	onMount(() => {
		get_teams();
	});

	async function get_teams(refresh) {
		try {
			teams = [];
			err = "";

			let url = "api/config/get_teams";

			if (refresh) {
				url = url + "?refresh=y";
			}

			const response = await fetch(url);

			if (response.ok) {
				teams = await response.json();
			} else {
				throw new Error(await response.text());
			}
		} catch (error) {
			err = error.message;
		}
	}

	async function set_teams() {
		const data = teams.map((team) => ({
			slug: team.slug,
			review_enabled: team.review_enabled,
			review_order: team.review_order,
		}));

		try {
			err = "";
			result = "";

			const response = await fetch("api/config/set_teams", {
				method: "POST",
				headers: {
					"Content-Type": "application/json",
				},
				body: JSON.stringify(data),
			});

			result = await response.text();
		} catch (error) {
			err = error.message;
		}
	}
</script>

<h2>Team Configuration</h2>

<Button color="green" onClick={() => get_teams(true)}
	>Hard refresh team list</Button>
<Button color="green" onClick={() => set_teams()}>Save Teams</Button>

{#if err !== ""}
	<p>
		{err}
	</p>
{:else if teams.length > 0}
	<p>{teams.length} teams found</p>
	<table>
		<thead>
			<th>Enable Team</th>

			<th>Review Order</th>
		</thead>
		<tbody>
			{#each teams as team}
				<tr>
					<td>
						<label>
							{team.name}
							<input
								type="checkbox"
								id={team.slug}
								name={team.slug}
								bind:checked={team.review_enabled} />
						</label>
					</td>
					<td>
						<input
							type="number"
							min="0"
							max={teams.length}
							bind:value={team.review_order}
							disabled={!team.review_enabled} />
					</td>
				</tr>
			{/each}
		</tbody>
	</table>
	{#if result !== ""}
		<p>
			{result}
		</p>
	{/if}
{:else}
	<p>No teams found</p>
{/if}

<style>
</style>
