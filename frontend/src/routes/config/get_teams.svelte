<script>
  import { onMount } from "svelte";

  let teams = [];

  onMount(() => {
    get_teams();
  });

  async function get_teams() {
    try {
      const response = await fetch("http://localhost:8080/config/get_teams");

      teams = await response.json();
    } catch (error) {
      console.error("Error fetching data from the backend:", error);
    }
  }
</script>

<section>
  <h2>Team Configuration</h2>
  <button on:click={() => get_teams()}>Refresh Team List</button>

  {#if teams.length > 0}
    <p>{teams.length} teams found</p>
    <table>
      <thead>
        <th>Team Name</th>
        <th>Enable Team</th>
        <th>Review Order</th>
      </thead>
      <tbody>
        {#each teams as team}
          <tr>
            <td>{team.name}</td>
            <td>
              <input
                type="checkbox"
                id={team.slug}
                name={team.slug}
                value={team.slug}
              />
            </td>
            <td>
              <p>1</p>
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  {:else}
    <p>No teams found</p>
  {/if}
</section>

<style>
</style>
