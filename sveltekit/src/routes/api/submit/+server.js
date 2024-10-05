import { json } from "@sveltejs/kit";

/**
 * Handles POST requests to create a new todo item.
 *
 * @param {Object} context - The request context.
 * @param {Request} context.request - The incoming request object.
 * @returns {Promise<Response>} The response object containing the result of the POST request.
 */
export async function POST({ request }) {
    let todo = { status: 0, description: "placeholder" };

    try {
        const data = await request.json();
        todo.description = data.description;

        const response = await fetch("http://go:8090/todos", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(todo),
        });

        if (!response.ok) {
            console.error(`Failed to fetch: ${response.statusText}`);
            return json(
                { error: `Failed to fetch: ${response.statusText}` },
                { status: 503 },
            );
        }

        const result = await response.json();
        return json(result);
    } catch (/** @type {any} */ error) {
        console.error(`Error processing request: ${error.message}`);
        return json(
            { error: `Error processing request: ${error.message}` },
            { status: 500 },
        );
    }
}
