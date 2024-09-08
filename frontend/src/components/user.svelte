<script>
  export let user;
</script>

<a href={user.html_url} target="_blank" class="container">
  <img src={user.avatar_url} alt="{user.login} avatar" />
  <div class="name-container">
    {#if user.name !== undefined}
      <span class="name">{user.name}</span>
      <span class="login">@{user.login}</span>
    {:else}
      <span class="big-login">@{user.login}</span>
    {/if}
  </div>
</a>

<style>
  a.container {
    margin: 12px;
    display: inline-flex;
    align-items: center;
    color: var(--text);
  }

  a.container:hover {
    color: var(--text-links);
  }

  img {
    width: 32px;
    height: 32px;
    border-radius: 50%;
    margin-right: 8px;
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
