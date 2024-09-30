<script>
  import { onMount } from "svelte";
  import User from "../../components/user.svelte";

  let members = [];
  let err = "";

  let team_members = {};
  let teams = [];

  onMount(() => {
    get_members();
  });

  async function get_members(refresh) {
    try {
      err = "";
      members = [];
      team_members = { none: [] };
      teams = [{ name: "none" }];

      let url = "http://localhost:8080/config/get_members";

      if (refresh) {
        url = url + "?refresh=y";
      }

      const response = await fetch(url);

      if (response.ok) {
        members = await response.json();
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
      } else {
        throw new Error(await response.text());
      }
    } catch (error) {
      err = error.message;
    }
  }
</script>

<h2>Member Configuration</h2>
<button on:click={() => get_members(true)}>Hard refresh member list</button>

{#if err !== ""}
  <p>
    {err}
  </p>
{:else if members.length > 0}
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
