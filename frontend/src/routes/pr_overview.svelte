<script>
  import { onMount } from "svelte";
  import Review from "../components/review.svelte";
  import PRState from "../components/pr_state.svelte";
  import User from "../components/user.svelte";
  import PRStats from "../components/pr_stats.svelte";

  let pr_list = [];
  let pr_stats = { "ready to merge": 0, "Changes Requested": 0 };

  onMount(() => {
    get_pr_list();
  });

  async function get_pr_list() {
    try {
      const response = await fetch("http://localhost:8080/get_pr_list");

      pr_list = await response.json();

      if (response.ok) {
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

<section>
  <h1>Github Pull request overview</h1>
  {#if pr_stats != undefined}
    {#each Object.entries(pr_stats) as [who, count]}
      <PRStats {who} {count} />
    {/each}
  {/if}

  {#if pr_list != undefined && pr_list.length > 0}
    number of pull requests {pr_list.length}
    <table>
      <tbody>
        {#each pr_list as pr}
          <tr>
            <td>
              <User user={pr.created_by} />
            </td>
            <td>
              <PRState state={pr.state} />
              #{pr.number}
            </td>
            <td>
              {pr.title}
              {#if pr.base.ref != "main"}
                <br />
                <span class="base">
                  {pr.base.ref}
                </span>
              {/if}
            </td>
            <td>
              {#if pr.awaiting != undefined}
                {pr.awaiting}
              {/if}
            </td>
            <td>
              {#if pr.review_overview !== undefined}
                {#each pr.review_overview as review}
                  {#if review.user !== undefined}
                    {review.user.login}
                  {:else}
                    {review.team.name}
                  {/if}
                  <Review state={review.state} /><br />
                {/each}
              {/if}
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  {:else}
    <p>No Pull Requests found</p>
  {/if}
</section>

<style>
  table {
    border-collapse: collapse;
    width: 100%;
    border-radius: 10%;
    border: 1px solid var(--border);
  }

  th,
  td {
    text-align: left;
    padding: 8px;
    border-top: 1px solid var(--border);
  }

  span.base {
    color: var(--yellow);
    font-weight: bold;
  }
</style>
