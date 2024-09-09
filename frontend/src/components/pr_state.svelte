<script>
  import Icon from "./icon.svelte";

  export let state;

  let width = "16px";
  let height = "16px";
</script>

{#if state === "open"}
  <Icon name="git-pull-request-16" color="green" {width} {height} />
{:else if state === "draft"}
  <Icon name="git-pull-request-draft-16" color="grey" {width} {height} />
{/if}
