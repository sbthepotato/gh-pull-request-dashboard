<script>
  let members = [];

  async function get_members() {
    try {
      const response = await fetch("http://localhost:8080/config/get_members");

      members = await response.json();
    } catch (error) {
      console.error("Error fetching data from the backend:", error);
    }
  }
</script>

<section>
  <h2>Connection test</h2>
  <button on:click={() => get_members}>Refresh contributor list</button>

  {#if members.length > 0}
    <p>{members.length}</p>
    <ul>
      {#each members as i}
        <li>{i.login}</li>
      {/each}
    </ul>
  {:else}
    <p>No members found</p>
  {/if}
</section>
