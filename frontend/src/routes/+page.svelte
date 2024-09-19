<script>
  import { onMount } from "svelte";

  import Button from "../components/button.svelte";
  import Icon from "../components/icon.svelte";
  import PRTable from "./pr_table.svelte";
  import PRAgg from "./pr_aggregation.svelte";

  let loading = false;
  let url = "http://localhost:8080/get_pr_list";
  let pr_list = [];
  let pr_stats = { total: 0, "ready to merge": 0, "Changes Requested": 0 };

  onMount(() => {
    get_pr_list();
  });

  async function get_pr_list(refresh) {
    try {
      loading = true;
      pr_list = [];
      pr_stats = { total: 0, "ready to merge": 0, "Changes Requested": 0 };

      if (refresh) {
        url = url + "?refresh=y";
      }

      const response = await fetch(url);

      pr_list = await response.json();

      if (response.ok) {
        pr_stats["total"] = pr_list.length;
        pr_list.forEach((pull) => {
          if (pull.awaiting === undefined) {
            pr_stats["ready to merge"] = pr_stats["ready to merge"] + 1 || 1;
          } else {
            pr_stats[pull.awaiting] = pr_stats[pull.awaiting] + 1 || 1;
          }
        });
      }
    } catch (error) {
      console.error("Error fetching data from the backend:", error);
    } finally {
      loading = false;
    }
  }
</script>

{#if loading}
  <div>
    <p>Loading PR list...</p>
    <Icon name="mark-github-24" color="rainbow" height="128px" width="128px" />
  </div>
{:else}
  <PRAgg {pr_stats} />
  <PRTable {pr_list} />
{/if}

<Button to="/config">Config</Button>
<Button onClick={() => get_pr_list(true)}>Hard Refresh PR List</Button>
