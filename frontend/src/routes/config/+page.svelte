<script>
  import Button from "../../components/button.svelte";
  import HelloGo from "./hello_go.svelte";
  import MemberConfig from "./member_config.svelte";
  import TeamConfig from "./team_config.svelte";
</script>

<Button action="change_page" page="/" text="Back to Home" />

<TeamConfig />
<MemberConfig />
<HelloGo />
