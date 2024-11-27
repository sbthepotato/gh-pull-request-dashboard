<script>
	import Button from "../../components/button.svelte";
	import HelloGo from "./hello_go.svelte";
	import TeamConfig from "./team_config.svelte";
	import MemberConfig from "./member_config.svelte";
</script>

<Button to="/">Back to home</Button>

<div class="container">
	<TeamConfig />
</div>
<div class="container">
	<MemberConfig />
</div>
<div class="container">
	<HelloGo />
</div>

<footer>
	made by sbthepotato | <a
		href="https://github.com/sbthepotato/gh-pull-request-dashboard"
		target="_blank">GitHub Repo</a>
</footer>

<style>
	div.container {
		margin: 8px 4px 48px;
	}

	footer {
		margin: 48px 4px 16px;
	}
</style>
