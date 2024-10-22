<script>
	import { onDestroy, onMount } from "svelte";
	import { page } from "$app/stores";

	import Button from "../components/button.svelte";
	import Loading from "../components/loading.svelte";
	import PRTable from "./pr_table.svelte";
	import PRAgg from "./pr_aggregation.svelte";

	let loading = false;
	let auto_reload = false;
	let reload_interval;

	let url = "api/dashboard/get_pr_list";
	let err = "";
	let result = {};
	let created_by = "";
	let pr_list = {};

	onMount(() => {
		get_pr_list();
	});

	async function get_pr_list(refresh) {
		try {
			loading = true;
			err = "";
			pr_list = {};

			if (refresh) {
				url = url + "?refresh=y";
			}

			const response = await fetch(url);

			if (response.ok) {
				result = await response.json();

				pr_list = result.pull_requests;
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

	$: (created_by = $page.url.searchParams.get("created_by")),
		get_filter(created_by);

	$: result, get_filter($page.url.searchParams.get("created_by"));

	onDestroy(() => {
		clearInterval(reload_interval);
	});

	function get_filter(name) {
		if (name !== null && result.pull_requests !== undefined) {
			pr_list = result.pull_requests.filter(
				(pr) => pr.created_by.login == name,
			);
		} else {
			pr_list = result.pull_requests;
		}
	}
</script>

<section class="pr-table">
	{#if err !== ""}
		{err}
	{:else if loading}
		<Loading text="Loading PR list..." />
	{:else}
		<PRAgg {pr_list} review_teams={result.review_teams} />
		<PRTable {pr_list} />
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
