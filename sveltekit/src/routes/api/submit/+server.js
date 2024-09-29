import { json } from '@sveltejs/kit';

export async function POST({ request }) {

  let todo = { "status": 0, "description" : "placeholder"}

  const data = await request.json();
  todo.description = data.description;

  const response = await fetch('http://go:8090/todos', {
      method: 'POST',
      headers: {
          'Content-Type': 'application/json'
      },
      body: JSON.stringify(todo)
  });

  if (!response.ok) {
      console.error(`Failed to fetch: ${response.statusText}`);
      return json({ error: `Failed to fetch: ${response.statusText}` }, { status: 503 })
  }

  const result = await response.json();
  return json(result);
}
