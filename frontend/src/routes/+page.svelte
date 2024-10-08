<script>
	import { onDestroy, onMount } from "svelte";

	import Button from "../components/button.svelte";
	import Icon from "../components/icon.svelte";
	import PRTable from "./pr_table.svelte";
	import PRAgg from "./pr_aggregation.svelte";

	let loading = false;
	let auto_reload = false;
	let reload_interval;

	let url = "api/dashboard/get_pr_list";
	let err = "";
	let pr_list = {};
	let pr_stats = {};

	onMount(() => {
		get_pr_list();
	});

	async function get_pr_list(refresh) {
		try {
			loading = true;
			err = "";
			pr_list = {};
			pr_stats = { total: 0, "ready to merge": 0, "Changes Requested": 0 };

			if (refresh) {
				url = url + "?refresh=y";
			}

			const response = await fetch(url);

			if (response.ok) {
				pr_list = await response.json();
				pr_stats["total"] = pr_list.pull_requests.length;

				// make sure the review teams are in the correct order always
				if (pr_list.review_teams !== undefined) {
					pr_list.review_teams.forEach((team) => {
						pr_stats[team.name] = 0;
					});
				} else {
					pr_stats[review] = 0;
				}

				pr_list.pull_requests.forEach((pull) => {
					if (pull.awaiting === "APPROVED") {
						pr_stats["ready to merge"] = pr_stats["ready to merge"] + 1 || 1;
					} else if (pull.awaiting === undefined) {
						pr_stats["missing status"] = pr_stats["missing status"] + 1 || 1;
					} else {
						pr_stats[pull.awaiting] = pr_stats[pull.awaiting] + 1 || 1;
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

	$: if (auto_reload) {
		reload_interval = setInterval(get_pr_list, 600000);
	} else {
		clearInterval(reload_interval);
	}

	onDestroy(() => {
		clearInterval(reload_interval);
	});
</script>

<section class="pr-table">
	{#if err !== ""}
		{err}
	{:else if loading}
		<div>
			<p>Loading PR list...</p>
			<Icon
				name="mark-github-24"
				color="rainbow"
				height="128px"
				width="128px" />
		</div>
	{:else}
		<PRAgg {pr_stats} />
		<PRTable pr_list={pr_list.pull_requests} />
	{/if}
</section>

<section class="buttons">
	<Button color="grey" to="/config">Config</Button>
	<Button color="green" onClick={() => get_pr_list(true)}
		>Hard Refresh PR List</Button>
	<label>
		<input type="checkbox" bind:checked={auto_reload} />
		Auto refresh
	</label>
</section>

<style>
	section.pr-table {
		margin-bottom: 32px;
	}
</style>
