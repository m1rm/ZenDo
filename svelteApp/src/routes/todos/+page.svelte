<script>
    import { onMount } from 'svelte'
    import { todoData } from '../../stores/todos'

    onMount(async () => {
        fetch("http://localhost:8090/todos")
        .then(response => response.json())
        .then(data => {todoData.set(data)})
        .catch(error => {
            console.log(error)
            return []
        })
    })

    /**
     * @param {boolean} success
     * @param {string} action
     * @param {Response} response
     */
    async function builtConfirmationMessage(success, action, response) {
        if (!success) {
            confirmationMessage = `Failed to ${action} todo: ${response.statusText}`
            showConfirmation = true
            console.log(await response.json())
        } else {
            confirmationMessage = `${action} todo: success.`
                textInput = ''; // Reset the form input
                showConfirmation = true;

                // Clear the confirmation message after 3 seconds
                setTimeout(() => {
                    showConfirmation = false
                }, 2000);
        }
    }

    let textInput = '';
    let confirmationMessage = '';
    let showConfirmation = false;
    /**
     * @param {Event} event
     */
    async function handleSubmit(event) {
        event.preventDefault();

        try {
            const response = await fetch('/api/submit', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ text: textInput })
            })

            builtConfirmationMessage(response.ok, 'add', response)
            if(!response.ok) {
                throw new Error('Failed to add the todo')
            } else {
                // @todo: handle success -> update UI
            }
        } catch (error) {
            console.error('Error:', error)
        }
    }

    /**
     * @param {number} id
     */
    async function deleteTodo(id) {
    try {
        const response = await fetch(`http://localhost:8090/todos/${id}`, {
            method: 'DELETE'
        })

        builtConfirmationMessage(response.ok, 'delete', response)
        if(!response.ok) {
            throw new Error('Failed to delete the todo')
        } else {
            // @todo: handle success -> update UI
        }
    } catch (error) {
      console.error('Error:', error)
    }
  }

</script>

<span class="h1 col-12">ToDos</span>

{#if showConfirmation}
<div style="z-index: 11" class="position-fixed top-50 start-50 translate-middle">
    <div class="toast show bg-info" role="alert" aria-live="assertive" aria-atomic="true">
        <div class="toast-header">
            <strong class="me-auto">Info</strong>
            <button type="button" class="btn-close" data-bs-dismiss="toast" aria-label="Close"></button>
        </div>
        <div class="toast-body">
            {confirmationMessage}
        </div>
    </div>
</div>
{/if}

<div class="row m-2 gap-2">
    <form on:submit={handleSubmit}>
        <div class="form-group mb-1 col-6">
            <label for="textInput">Enter Text</label>
            <input type="text" class="form-control" id="textInput" bind:value={textInput} placeholder="Enter text">
        </div>
        <button type="submit" class="btn btn-success">Submit</button>
    </form>

    {#each $todoData as todo}
        <div class="col-12 col-sm-6 border d-flex flex-column">
            <span class="mt-2 mb-2">
            { todo.text }
            </span>

            <div class="d-flex justify-content-between align-items-end gap-2 mb-2">
                <label>
                    <input class="form-check-input flex-shrink-0" type="checkbox" checked="{todo.status === 1}">
                    { todo.status === 1 ? 'done' : 'open' }
                </label>
                <div>
                    <button class="btn btn-sm btn-outline-secondary" type="button">Edit</button>
                    <button 
                        class="btn btn-sm btn-primary btn-danger"
                        type="button"
                        on:click={() => deleteTodo(todo.id)}>
                        Delete
                    </button>
                </div>
            </div>
        </div>
    {/each}
</div>
