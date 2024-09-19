<script>
  import Review from "../components/review.svelte";
  import PRState from "../components/pr_state.svelte";
  import User from "../components/user.svelte";
  import Icon from "../components/icon.svelte";

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

  function getTextLuminance(hexColor) {
    // Remove the "#" if it exists
    const color = hexColor.replace("#", "");

    // Convert hex to RGB
    const r = parseInt(color.substring(0, 2), 16);
    const g = parseInt(color.substring(2, 4), 16);
    const b = parseInt(color.substring(4, 6), 16);

    // Calculate the luminance of the color
    const luminance = (0.299 * r + 0.587 * g + 0.114 * b) / 255;

    // Return white for dark colors and black for light colors
    return luminance > 0.5 ? "#0d1117" : "#f0f6fc";
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
            {#if pr.labels != undefined}
              {#each pr.labels as label}
                <span
                  class="tag"
                  style="background-color:#{label.color}; color: {getTextLuminance(
                    label.color
                  )}"
                >
                  {label.name}
                </span>
                &nbsp;
              {/each}
            {/if}
            <br />
            <span class="under-text">
              <span>
                <a href={pr.html_url} class="pr_url">
                  #{pr.number}
                </a>
              </span>
              <span>
                Created {convert_date(pr.created_at)}
              </span>
              <span>
                Modified {convert_date(pr.updated_at)}
              </span>
              <span>
                <Icon
                  name="file-code-16"
                  height="12px"
                  width="12px"
                  color="grey"
                />
                {pr.changed_files}
              </span>
              <span class="additions"> ++{pr.additions}</span>
              <span class="deletions">--{pr.deletions}</span>
              <span
                ><Icon
                  name="comment-discussion-16"
                  height="12px"
                  width="12px"
                  color="grey"
                />
                {pr.comments}
              </span>
              {#if pr.base.ref != "main"}
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

  span.tag {
    border-radius: 400px;
    padding: 2px 8px;
    font-size: small;
    font-weight: bold;
  }

  span.under-text {
    font-size: small;
    color: var(--text-alt);
    margin-left: 28px;
  }
  span.under-text > span {
    padding: 0px 4px;
  }

  span.additions {
    color: var(--green);
  }

  span.deletions {
    color: var(--red);
  }

  span.base {
    color: var(--yellow);
    font-weight: bold;
  }

  a.pr_url:hover {
    text-decoration: underline;
  }
</style>
