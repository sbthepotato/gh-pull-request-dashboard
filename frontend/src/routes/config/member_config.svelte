<script>
  import { onMount } from "svelte";
  import User from "../../components/user.svelte";

  let members = [];

  let team_members = { none: [] };
  let teams = [{ name: "none" }];

  onMount(() => {
    get_members();
  });

  async function get_members(refresh) {
    try {
      let url = "http://localhost:8080/config/get_members";

      if (refresh) {
        url = url + "?refresh=y";
      }

      const response = await fetch(url);

      members = await response.json();

      if (response.ok) {
        members.forEach((user) => {
          if (user.team === undefined) {
            team_members["none"].push(user);
          } else if (user.team.name in team_members) {
            team_members[user.team.name].push(user);
          } else {
            team_members[user.team.name] = [user];
            teams.push(user.team);
          }
        });
      }
    } catch (error) {
      console.error("Error fetching data from the backend:", error);
    }
  }
</script>

<h2>Member Configuration</h2>
<button on:click={() => get_members(true)}>Hard refresh member list</button>

{#if members.length > 0}
  <p>{members.length} members found</p>
  {#each teams as team}
    <div class="team-container">
      <h2>{team.name}</h2>

      {#each team_members[team.name] as member}
        <User user={member} />
      {/each}
    </div>
  {/each}
{:else}
  <p>No members found</p>
{/if}

<style>
  div.team-container {
    border: 1px solid var(--border-blue);
    margin: 8px;
    border-radius: 5px;
  }
</style>
