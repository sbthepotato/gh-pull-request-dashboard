<script>
  import Review from "../components/review.svelte";
  import PRState from "../components/pr_state.svelte";
  import User from "../components/user.svelte";

  export let pr_list = [];

  let date_options = {
    month: "short",
    day: "numeric",
    hour: "numeric",
    minute: "numeric",
    hour12: false,
  };

  function convert_date(date_str) {
    let date = new Date(date_str);
    return date.toLocaleString("en-us", date_options);
  }
</script>

{#if pr_list != undefined && pr_list.length > 0}
  <table>
    <tbody>
      {#each pr_list as pr}
        <tr>
          <td>
            <User user={pr.created_by} />
          </td>

          <td>
            <PRState state={pr.state} />
            &nbsp;
            {pr.title}
            <br />
            <span class="under-text">
              <a href={pr.html_url} class="pr_url">
                #{pr.number}
              </a>
              - Created {convert_date(pr.created_at)}
              - Modified {convert_date(pr.updated_at)}
              {#if pr.base.ref != "main"}
                -
                <span class="base">
                  {pr.base.ref}
                </span>
              {/if}
            </span>
          </td>
          <td>
            {#if pr.awaiting != undefined}
              {pr.awaiting}
            {/if}
          </td>
          <td>
            {#if pr.review_overview !== undefined}
              {#each pr.review_overview as review}
                {#if review.user !== undefined}
                  <User user={review.user} size="xs" />
                {:else}
                  {review.team.name}
                {/if}
                <Review state={review.state} /><br />
              {/each}
            {/if}
          </td>
        </tr>
      {/each}
    </tbody>
  </table>
{:else}
  <p>No Pull Requests found</p>
{/if}

<style>
  table {
    border-collapse: collapse;
    width: 100%;
    border-radius: 10%;
    border: 1px solid var(--border);
  }

  td {
    text-align: left;
    padding: 8px;
    border-top: 1px solid var(--border);
  }

  span.under-text {
    font-size: small;
    color: var(--text-alt);
    margin-left: 28px;
  }

  span.base {
    color: var(--yellow);
    font-weight: bold;
  }

  a.pr_url:hover {
    text-decoration: underline;
  }
</style>
