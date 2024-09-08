<script>
  import Icon from "./icon.svelte";

  export let state;

  let width = "24px";
  let height = "24px";
</script>

{#if state === "open"}
  <Icon name="git-pull-request" color="green" {width} {height} />
{:else if state === "draft"}
  <Icon name="git-pull-request-draft" color="grey" {width} {height} />
{/if}
