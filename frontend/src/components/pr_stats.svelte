<script>
  import Icon from "./icon.svelte";

  export let who = "";
  export let count = 0;

  let height = "16px";
  let width = "16px";
</script>

<div class="container">
  {#if who === "ready to merge"}
    <p>Ready to Merge</p>
    <Icon name="check-16" color="green" {height} {width} />
  {:else if who === "Changes Requested"}
    <p>Changes Requested</p>
    <Icon name="file-diff-16" color="red" {height} {width} />
  {:else if who === "total"}
    <p>Total Pull Requests</p>
    <Icon name="git-pull-request-16" color="green" {height} {width} />
  {:else}
    <p>Waiting on {who}</p>
    <Icon name="dot-fill-16" color="yellow" {height} {width} />
  {/if}

  <p class="count">{count}</p>
</div>

<style>
  div.container {
    display: inline-block;
    padding: 8px;
    margin: 8px;
    border-radius: 8px;
    border: 1px solid var(--border);
  }

  p {
    margin: 4px 0;
  }

  p.count {
    font-weight: bold;
    font-size: x-large;
  }
</style>
