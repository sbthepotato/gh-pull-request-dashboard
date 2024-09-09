<script>
  import Icon from "./icon.svelte";

  export let state;

  let height = "16px";
  let width = "16px";
</script>

{#if state === "APPROVED"}
  <Icon name="check-16" color="green" {height} {width} />
{:else if state === "CHANGES_REQUESTED"}
  <Icon name="file-diff-16" color="red" {height} {width} />
{:else if state === "COMMENTED"}
  <Icon name="comment-16" color="yellow" {height} {width} />
{:else if state === "REVIEW_REQUESTED"}
  <Icon name="dot-fill-16" color="yellow" {height} {width} />
{/if}
