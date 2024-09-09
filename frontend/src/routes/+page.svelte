<script>
  import { onMount } from "svelte";

  import Button from "../components/button.svelte";
  import PRTable from "./pr_table.svelte";
  import PRAgg from "./pr_aggregation.svelte";

  let url = "http://localhost:8080/get_pr_list";
  let pr_list = [];
  let pr_stats = { total: 0, "ready to merge": 0, "Changes Requested": 0 };

  onMount(() => {
    get_pr_list();
  });

  async function get_pr_list(refresh) {
    try {
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
    }
  }
</script>

<PRAgg {pr_stats} />
<PRTable {pr_list} />

<Button to="/config">Config</Button>
<Button onClick={() => get_pr_list(true)}>Hard Refresh PR List</Button>
