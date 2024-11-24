<script>
	import { onDestroy, onMount } from "svelte";
	import { page } from "$app/stores";
	import {
		set_many_url_params,
		set_url_param,
		string_to_bool,
	} from "$lib/index.js";

	import Button from "../components/button.svelte";
	import Checkbox from "../components/checkbox.svelte";
	import Searchbar from "../components/searchbar.svelte";
	import Loading from "../components/loading.svelte";
	import PRTable from "./pr_table.svelte";
	import PRAgg from "./pr_aggregation.svelte";

	let url = "api/dashboard/get_pr_list";
	let err = "";
	let result = {};
	let pr_list = {};

	let loading = false;

	let checkboxes = {
		auto_reload: false,
		show_search: false,
	};
	let reload_interval;

	let created_by_filter = "";
	let created_by_filter_user = {};
	let search_query = "";

	onMount(() => {
		get_pr_list();

		created_by_filter = $page.url.searchParams.get("created_by");

		checkboxes.auto_reload = string_to_bool(
			$page.url.searchParams.get("auto_reload"),
			false,
		);

		checkboxes.show_search = string_to_bool(
			$page.url.searchParams.get("show_search"),
			false,
		);
	});

	onDestroy(() => {
		clearInterval(reload_interval);
	});

	async function get_pr_list(refresh) {
		try {
			loading = true;
			err = "";
			result = {};
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

	function handle_checkbox_change(event) {
		const { id, checked } = event.detail;
		checkboxes = { ...checkboxes, [id]: checked };

		switch (id) {
			case "show_search":
				if (checked) {
					set_url_param("show_search", "y");
					get_filter();
				} else {
					set_url_param("show_search");
					get_filter();
				}
				break;
			case "auto_reload":
				if (checked) {
					set_url_param("auto_reload", "y");
					reload_interval = setInterval(get_pr_list, 600000);
				} else {
					set_url_param("auto_reload");
					clearInterval(reload_interval);
				}
				break;
		}
	}

	function handle_searchbar_change(event) {
		search_query = event.detail.value.toLowerCase();
		get_filter();
	}

	function handle_params() {
		created_by_filter = $page.url.searchParams.get("created_by");
		checkboxes.show_search = string_to_bool(
			$page.url.searchParams.get("show_search"),
			false,
		);
	}

	function get_filter() {
		if (
			(created_by_filter !== null || search_query !== "") &&
			result.pull_requests !== undefined
		) {
			pr_list = result.pull_requests.filter(
				(pr) =>
					(created_by_filter === null ||
						pr.created_by.login === created_by_filter ||
						pr.review_overview?.some(
							(review) =>
								review.user?.login === created_by_filter &&
								review.state === "REVIEW_REQUESTED",
						) ||
						(pr.unassigned === true &&
							pr.awaiting === created_by_filter_user.team?.name)) &&
					(pr.title.toLowerCase().includes(search_query) ||
						pr.awaiting?.toLowerCase().includes(search_query) ||
						pr.created_by.login.toLowerCase().includes(search_query) ||
						pr.created_by.name.toLowerCase().includes(search_query) ||
						pr.base.label.toLowerCase().includes(search_query) ||
						pr.number.toString().includes(search_query) ||
						pr.review_overview?.some(
							(review) =>
								review.state === "REVIEW_REQUESTED" &&
								(review.user?.login.toLowerCase().includes(search_query) ||
									review.user?.name.toLowerCase().includes(search_query)),
						) ||
						pr.labels?.some((label) =>
							label.name.toLowerCase().includes(search_query),
						)),
			);

			if (created_by_filter !== null) {
			}
		} else {
			pr_list = result.pull_requests;
		}
	}

	function get_current_user() {
		if (created_by_filter == null) {
			created_by_filter_user = {};
		} else {
			if (result.users !== undefined) {
				result.users.forEach((user) => {
					if (user?.login === created_by_filter) {
						created_by_filter_user = user;
						return true;
					}
				});
			}
		}
	}

	function clear_filters() {
		set_many_url_params({ created_by: null });
		created_by_filter = null;
		search_query = "";
		get_filter();
	}

	$: $page.url.search, handle_params();
	$: result, get_current_user(), get_filter();
	$: created_by_filter, get_current_user(), get_filter();
</script>

<section class="pr-table">
	{#if err !== ""}
		{err}
	{:else if loading}
		<Loading text="Loading PR list..." />
	{:else}
		<PRAgg {pr_list} review_teams={result.review_teams} />
		{#if checkboxes.show_search}
			<Searchbar
				value={search_query}
				placeholder="Search Pull Requests..."
				on:change={handle_searchbar_change}
				on:input={handle_searchbar_change} />
		{/if}
		{#if created_by_filter === null}
			<PRTable {pr_list} />
		{:else if created_by_filter !== null}
			<PRTable
				title="Created by {created_by_filter}"
				pr_list={pr_list?.filter(
					(pr) => pr.created_by.login === created_by_filter,
				)} />

			<PRTable
				title="{created_by_filter} requested reviewer"
				pr_list={pr_list?.filter((pr) =>
					pr.review_overview?.some(
						(review) =>
							review.user?.login === created_by_filter &&
							review.state === "REVIEW_REQUESTED",
					),
				)} />

			{#if created_by_filter_user.team}
				<PRTable
					title="Waiting on {created_by_filter_user.team
						.name} - Not assigned to anyone else"
					pr_list={pr_list?.filter(
						(pr) =>
							pr.unassigned === true &&
							pr.awaiting === created_by_filter_user.team.name,
					)} />
			{/if}
		{/if}
	{/if}
</section>

<section class="buttons">
	<Button color="grey" to="/config">Config</Button>
	<Button color="green" on_click={() => get_pr_list(true)}>
		Hard Refresh PR List
	</Button>
	<Checkbox
		id="auto_reload"
		checked={checkboxes.auto_reload}
		on:change={handle_checkbox_change}>Auto Refresh</Checkbox>
	<Checkbox
		id="show_search"
		checked={checkboxes.show_search}
		on:change={handle_checkbox_change}>Show Search</Checkbox>
	{#if created_by_filter !== null || search_query !== ""}
		<Button color="blue" on_click={() => clear_filters()}>Clear Filters</Button>
	{/if}
</section>

<style>
	section.pr-table {
		margin-bottom: 32px;
	}
</style>
