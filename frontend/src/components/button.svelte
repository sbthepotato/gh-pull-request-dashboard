<script>
  import { goto } from "$app/navigation";

  export let type = "button";
  export let onClick = () => {};
  export let to = null;

  function handleClick(event) {
    if (to) {
      goto(to);
    } else {
      onClick(event);
    }
  }
</script>

<button {type} on:click={handleClick}>
  <slot></slot>
</button>

<style>
  button {
    display: inline;
  }
</style>
