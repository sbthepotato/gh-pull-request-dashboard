// src/hooks.server.js
import { redirect } from "@sveltejs/kit";

export async function handle({ event, resolve }) {
	const url = new URL(event.request.url);

	if (!url.pathname.endsWith("/")) {
		url.pathname += "/";
		throw redirect(308, url.toString());
	}

	return resolve(event);
}
