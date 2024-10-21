<script>
	import Review from "../components/review.svelte";
	import PrAwaiting from "../components/pr_awaiting.svelte";
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
					<td class="created_by">
						{#if pr.created_by !== undefined}
							<User user={pr.created_by} action="filter" />
						{:else}
							<User user={pr.user} />
						{/if}
					</td>

					<td class="title">
						<span class="pr-title">{pr.title}</span>
						{#if pr.labels != undefined}
							<span class="tags">
								{#each pr.labels as label}
									&nbsp;
									<span
										class="tag"
										style="background-color:#{label.color}; color: {getTextLuminance(
											label.color,
										)}">
										{label.name}
									</span>
								{/each}
							</span>
						{/if}
						<br />
						<span class="under-text">
							<span>
								<a href={pr.html_url} target="_blank" class="pr_url">
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
									color="grey" />
								{pr.changed_files}
							</span>
							<span>
								<span class="additions"> ++{pr.additions}</span>
								<span class="deletions">--{pr.deletions}</span>
							</span>
							<span
								><Icon
									name="comment-discussion-16"
									height="12px"
									width="12px"
									&nbsp;
									color="grey" />
								{pr.comments}
							</span>
							{#if pr.base.ref != "main"}
								<span class="base">
									{pr.base.ref}
								</span>
							{/if}
						</span>
					</td>
					<td class="awaiting">
						<PrAwaiting awaiting={pr.awaiting} />
					</td>
					<td class="review_overview">
						{#if pr.review_overview !== undefined}
							{#each pr.review_overview as review}
								{#if review.user !== undefined}
									<User user={review.user} size="xs" />
								{:else if review.team !== undefined}
									{review.team.name}
								{:else}
									UNKNOWN
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
		width: 100%;
		border-spacing: 0;
		border-collapse: separate;
	}

	table tr td {
		border-top: 1px solid var(--border);
	}

	table td {
		text-align: left;
		padding: 8px;
	}

	table tr:last-child td {
		border-bottom: 1px solid var(--border);
	}

	table td:first-child {
		border-left: 1px solid var(--border);
	}

	table td:last-child {
		border-right: 1px solid var(--border);
	}

	table tr:first-child td:first-child {
		border-top-left-radius: 8px;
	}

	table tr:first-child td:last-child {
		border-top-right-radius: 8px;
	}

	table tr:last-child td:first-child {
		border-bottom-left-radius: 8px;
	}

	table tr:last-child td:last-child {
		border-bottom-right-radius: 8px;
	}

	span.tags {
		margin-left: 12px;
	}

	span.tag {
		border-radius: 400px;
		padding: 2px 8px;
		font-size: small;
		font-weight: bold;
		white-space: nowrap;
	}

	span.under-text {
		font-size: small;
		color: var(--text-alt);
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

	td.awaiting {
		min-width: 100px;
	}
</style>
