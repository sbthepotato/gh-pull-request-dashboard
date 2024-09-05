<script>
  import Icon from "./icon.svelte";

  export let state;

  let height = "16px";
  let width = "16px";
</script>

{#if state === "APPROVED"}
  <Icon name="check" color="green" {height} {width} />
{:else if state === "CHANGES_REQUESTED"}
  <Icon name="file-diff" color="red" {height} {width} />
{:else if state === "COMMENTED"}
  <Icon name="comment" color="yellow" {height} {width} />
{:else if state === "REQUESTED"}
  <Icon name="dot-fill" color="yellow" {height} {width} />
{/if}
