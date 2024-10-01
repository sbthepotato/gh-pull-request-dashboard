<script>
  let answer = "";
  let err = "";

  async function hello_go() {
    try {
      answer = "";
      err = "";

      const response = await fetch("/api/config/hello_go");

      if (response.ok) {
        answer = await response.text();
      } else {
        throw new Error(await response.text());
      }
    } catch (error) {
      err = error.message;
    }
  }
</script>

<h2>Connection test</h2>
<button on:click={hello_go}>Say hello to the backend</button>

<p>
  {#if err}
    <br />
    <span class="bad">{err}</span>
  {:else if answer}
    <br />
    <span class="good">{answer}</span>
  {/if}
</p>

<style>
  span.good {
    color: var(--green);
  }

  span.bad {
    color: var(--red);
  }
</style>
