<script>
    import { onMount } from 'svelte'
    import { page } from '$app/stores'
    import { todoData, selectedTodo } from '../../../stores/todos'

    onMount(async () => {
        const id = $page.params.id
        try {
            // Fetch all todos
            const response = await fetch('http://localhost:8090/todos')
            const data = await response.json()
            todoData.set(data)

            const todo = data.find(todo => todo.id === id)
            selectedTodo.set(todo)
        } catch (error) {
            console.log(error)
        }
    });

    /**
     * Handles form submission to update a todo item.
     * @param {Event} event
     */
    const handleSubmit = async (event) => {
        event.preventDefault();
        const id = $page.params.id;
        const formData = new FormData(event.target);
        const updatedTodo = {
            title: formData.get('title'),
            description: formData.get('description')
        };

        try {
            const response = await fetch(`http://localhost:8090/todos/${id}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(updatedTodo)
            });
            const data = await response.json();
            selectedTodo.set(data);
        } catch (error) {
            console.log(error);
        }
    }
</script>

<div>
    {#if $selectedTodo}
        <h1>#{ $selectedTodo.id }</h1>
        <p>{ $selectedTodo.description }</p>

        <!-- Edit form -->
        <form on:submit={handleSubmit}>
            <label for="description">Description:</label>
            <textarea id="description" name="description" required>{$selectedTodo.description}</textarea>

            <label for="status">Done:</label>
            <input type="checkbox" id="status" name="status" value="{$selectedTodo.status}">

            <button type="submit">Save</button>
        </form>
    {/if}
</div>
