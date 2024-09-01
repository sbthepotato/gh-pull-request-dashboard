<script>
  let contributors = [];

  async function listContributors() {
    try {
      const response = await fetch(
        "http://localhost:8080/config/get_contributors"
      );

      const data = await response.json();

      contributors = data.message;

      console.log(contributors);
    } catch (error) {
      console.error("Error fetching data from the backend:", error);
      answer = "Error connecting to the backend";
    }
  }
</script>

<section>
  <h2>Connection test</h2>
  <button on:click={listContributors}>Refresh contributor list</button>

  {#if contributors.length > 0}
    <p>{contributors.length}</p>
    <!--   <ul>
      {#each contributors as i}
        <li>{i.login}</li>
      {/each}
    </ul> -->
  {:else}
    <p>No contributors found</p>
  {/if}
</section>
