<script>
	import { onMount } from "svelte";
	import User from "../../components/user.svelte";
	import Button from "../../components/button.svelte";
	import Loading from "../../components/loading.svelte";

	let err = "";
	let result = [];
	let loading = false;

	let team_members = {};
	let teams = [];

	onMount(() => {
		get_members();
	});

	async function get_members(refresh) {
		try {
			loading = true;
			err = "";
			result = [];
			team_members = { none: [] };
			teams = [{ name: "none" }];

			let url = "api/config/get_members";

			if (refresh) {
				url = url + "?refresh=y";
			}

			const response = await fetch(url);

			if (response.ok) {
				result = await response.json();

				result.forEach((user) => {
					if (user.team === undefined) {
						team_members["none"].push(user);
					} else if (user.team.name in team_members) {
						team_members[user.team.name].push(user);
					} else {
						team_members[user.team.name] = [user];
						teams.push(user.team);
					}
				});
			} else {
				throw new Error(await response.text());
			}
		} catch (error) {
			err = error.message;
		} finally {
			loading = false;
		}
	}
</script>

<h2>Member Configuration</h2>
<Button color="green" on_click={() => get_members(true)}
	>Hard refresh member list</Button>

{#if err !== ""}
	<p>
		{err}
	</p>
{:else if loading}
	<Loading text="Loading Members..." size="64px" />
{:else if result.length > 0}
	<p>{result.length} members found</p>
	{#each teams as team}
		<div class="team-container">
			<h2>{team.name}</h2>

			{#each team_members[team.name] as member}
				<User user={member} />
			{/each}
		</div>
	{/each}
{:else}
	<p>No members found</p>
{/if}

<style>
	div.team-container {
		border: 1px solid var(--border-blue);
		margin: 8px;
		border-radius: 5px;
	}
</style>
