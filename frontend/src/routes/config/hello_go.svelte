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

<section id="hello_go">
  <h2>Connection test</h2>
  <button on:click={hellogo}>Say hello to the backend</button>

  {#if answer !== undefined}
    <p>
      Backend says:
      <br />
      <span>{answer}</span>
    </p>
  {/if}
</section>

<style>
  section {
    border: 1px solid black;
    display: inline-block;
    padding: 4px;
  }

  span {
    color: green;
  }
</style>
