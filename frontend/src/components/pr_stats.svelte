<script>
  export let who = "";
  export let count = 0;
</script>

<div class="container">
  {#if who === "ready to merge"}
    <p>Ready to Merge</p>
  {:else if who === "Changes Requested"}
    <p>Changes Requested</p>
  {:else}
    <p>Waiting on {who}</p>
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

  p.count {
    font-weight: bold;
    font-size: x-large;
    margin-bottom: 4px;
  }
</style>
