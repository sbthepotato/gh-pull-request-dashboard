<script>
  let answer;

  async function hellogo() {
    try {
      const response = await fetch("http://localhost:8080/config/hello_go");

      const data = await response;

      console.log(data);

      answer = data;
    } catch (error) {
      console.error("Error fetching data from the backend:", error);
      answer = "Error connecting to the backend";
    }
  }
</script>

<h2>Connection test</h2>
<button on:click={hellogo}>Say hello to the backend</button>

{#if answer !== undefined}
  <p>
    Backend says:
    <br />
    <span>{answer}</span>
  </p>
{/if}

<style>
  span {
    color: green;
  }
</style>
