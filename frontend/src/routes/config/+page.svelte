<script>
	import Button from "../../components/button.svelte";
	import HelloGo from "./hello_go.svelte";
	import TeamConfig from "./team_config.svelte";
	import MemberConfig from "./member_config.svelte";
</script>

<Button to="/">Back to home</Button>

<TeamConfig />
<MemberConfig />
<HelloGo />

<footer>
	made by sbthepotato | <a
		href="https://github.com/sbthepotato/gh-pull-request-dashboard"
		target="_blank">GitHub Repo</a>
</footer>

<style>
	footer {
		margin-top: 32px;
	}
</style>
