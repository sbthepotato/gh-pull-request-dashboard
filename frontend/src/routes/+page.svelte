<script>
  let pr_list = [];

  async function listContributors() {
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

  <button on:click={listContributors}>Refresh PR list</button>

  {#if pr_list.length > 0}
    <p>{pr_list.length}</p>
    <table>
      <thead>
        <tr>
          <th>Created By</th>
          <th>State</th>
          <th>Title</th>
        </tr>
      </thead>
      <tbody>
        {#each pr_list as pr}
          <tr>
            <td
              ><img
                src={pr.user.avatar_url}
                alt="github avatar of {pr.user.login}"
              />
              <a href={pr.user.html_url}>
                {pr.user.login}
              </a>
            </td>
            <td>
              {#if pr.draft}
                draft
              {:else}
                {pr.state}
              {/if}
            </td>
            <td>{pr.title}</td>
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

  td img {
    max-width: 50px;
  }

  th,
  td {
    text-align: left;
    padding: 8px;
  }

  tr:nth-child(odd) {
    background-color: lightgray;
  }
</style>
