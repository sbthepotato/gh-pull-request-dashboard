<script>
	import PRStats from "../components/pr_stats.svelte";

	export let pr_list;
	let pr_stats = {};

	function aggregate_prs(pr_list) {
		if (pr_list.pull_requests !== undefined) {
			pr_stats = { total: 0, "ready to merge": 0, "Changes Requested": 0 };
			pr_stats["total"] = pr_list.pull_requests.length;

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
		}
	}

	$: pr_stats, aggregate_prs(pr_list);
</script>

{#if pr_stats != undefined}
	{#each Object.entries(pr_stats) as [who, count]}
		<PRStats {who} {count} />
	{/each}
{/if}
