<script>
  import Button from "../../components/button.svelte";
  import HelloGo from "./hello_go.svelte";
  import GetMembers from "./get_members.svelte";
  import TeamConfig from "./team_config.svelte";
</script>

<Button action="change_page" page="/" text="Back to Home" />
<TeamConfig />
<HelloGo />
<GetMembers />
