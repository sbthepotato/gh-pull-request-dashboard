<script>
  import { onMount } from "svelte";
  import Review from "../components/review.svelte";
  import State from "../components/state.svelte";

  let pr_list = [];

  onMount(() => {
    get_pr_list();
  });

  async function get_pr_list() {
    try {
      const response = await fetch("http://localhost:8080/get_pr_list");

      pr_list = await response.json();
    } catch (error) {
      console.error("Error fetching data from the backend:", error);
    }
  }
</script>

<section>
  <h1>Github Pull request overview</h1>

  {#if pr_list.length > 0}
    <table>
      <tbody>
        {#each pr_list as pr}
          <tr>
            <td>
              <a href={pr.user.html_url}>
                {pr.user.login}
              </a>
            </td>
            <td>
              <State state={pr.state} />
              #{pr.number}
            </td>
            <td>
              {pr.title}
            </td>
            <td>
              {#if pr.Review_Overview !== undefined}
                {#each pr.Review_Overview as review}
                  {review.User} <Review state={review.State} /><br />
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
</style>
