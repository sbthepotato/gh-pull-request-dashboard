<script>
  import Button from "../../components/button.svelte";
  import HelloGo from "./hello_go.svelte";
  import GetMembers from "./get_members.svelte";
  import GetTeams from "./get_teams.svelte";
</script>

<Button action="change_page" page="/" text="Back to Home" />
<GetTeams />
<HelloGo />
<GetMembers />
