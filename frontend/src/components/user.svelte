<script>
	export let user;
	export let size = "l";
</script>

<a href={user.html_url} target="_blank" class="container">
	{#if size !== "xs"}
		<img src={user.avatar_url} alt="{user.login} avatar" class={size} />
	{/if}
	<div class="name-container">
		{#if user.name !== undefined}
			<span class="name">{user.name}</span>
			{#if size !== "xs"}
				<span class="login">@{user.login}</span>
			{/if}
		{:else}
			<span class="big-login">@{user.login}</span>
		{/if}
	</div>
</a>

<style>
	a.container {
		display: inline-flex;
		align-items: center;
		color: var(--text);
		font-size: medium;
	}

	a.container:hover {
		color: var(--text-links);
	}

	a.container:hover span.name,
	a.container:hover span.big-login {
		text-decoration: underline;
	}

	img {
		width: 32px;
		height: 32px;
		border-radius: 50%;
		margin-right: 8px;
	}

	img.s {
		width: 20px;
		height: 20px;
	}

	div.name-container {
		display: flex;
		flex-direction: column;
	}

	div.name-container > span {
		display: block;
		text-align: left;
	}

	span.login {
		color: var(--text-alt);
	}
</style>
