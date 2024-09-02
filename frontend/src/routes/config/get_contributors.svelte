<script>
  let contributors = [];

  async function listContributors() {
    try {
      const response = await fetch(
        "http://localhost:8080/config/get_contributors"
      );

      contributors = await response.json();
    } catch (error) {
      console.error("Error fetching data from the backend:", error);
    }
  }
</script>

<section>
  <h2>Connection test</h2>
  <button on:click={listContributors}>Refresh contributor list</button>

  {#if contributors.length > 0}
    <p>{contributors.length}</p>
    <ul>
      {#each contributors as i}
        <li>{i.login}</li>
      {/each}
    </ul>
  {:else}
    <p>No contributors found</p>
  {/if}
</section>
