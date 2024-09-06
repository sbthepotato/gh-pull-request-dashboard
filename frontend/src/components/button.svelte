<script>
  import { goto } from "$app/navigation";
  export let text;
  export let action;
  export let page;

  function button_action() {
    if (action === "change_page") {
      goto(page);
    }
  }
</script>

<button on:click={() => button_action()}>
  {text}
</button>

<style>
  button {
    display: inline;
  }
</style>
