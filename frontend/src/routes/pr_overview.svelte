<script>
  let pr_list = [];
  let reviews = [];

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

  <button on:click={get_pr_list}>Refresh PR list</button>

  {#if pr_list.length > 0}
    <table>
      <thead>
        <tr>
          <th>Created By</th>
          <th>State</th>
          <th>Title</th>
          <th>Reviews</th>
        </tr>
      </thead>
      <tbody>
        {#each pr_list as pr}
          <tr>
            <td>
              <a href={pr.user.html_url}>
                {pr.user.login}
              </a>
            </td>
            <td>
              {#if pr.draft}
                draft
              {:else}
                {pr.state} - {pr.number}
              {/if}
            </td>
            <td>{pr.title}</td>
            <td>
              {#if pr.Review_Overview !== undefined}
                {#each pr.Review_Overview as review}
                  {review.User} - {review.State} <br />
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
  }

  th,
  td {
    text-align: left;
    padding: 8px;
  }

  tr:nth-child(odd) {
    background-color: #2b2b33;
  }
</style>
