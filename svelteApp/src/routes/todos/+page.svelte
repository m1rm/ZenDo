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

    let textInput = '';
    let confirmationMessage = '';
    let showConfirmation = false;
    /**
     * @param {Event} event
     */
    async function handleSubmit(event) {
        event.preventDefault();

        const response = await fetch('/api/submit', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ text: textInput })
        })

        if (!response.ok) {
            confirmationMessage = `Failed to add todo: ${response.statusText}`;
            showConfirmation = true;
            console.log(await response.json());
        } else {
            confirmationMessage = 'Todo added successfully!';
                textInput = ''; // Reset the form input
                showConfirmation = true;

                // Clear the confirmation message after 3 seconds
                setTimeout(() => {
                    showConfirmation = false;
                }, 3000);
            const result = await response.json();
            console.log(result);
        }
    }

</script>

<style>
    .fade {
        transition: opacity 0.5s ease-in-out;
    }
    .fade.show {
        opacity: 1;
    }
    .fade:not(.show) {
        opacity: 0;
    }
</style>


<span class="h1 col-12">ToDos</span>

<div class="row m-2 gap-2">
    <form on:submit={handleSubmit}>
        <div class="form-group mb-1 col-6">
            <label for="textInput">Enter Text</label>
            <input type="text" class="form-control" id="textInput" bind:value={textInput} placeholder="Enter text">
        </div>
        <button type="submit" class="btn btn-success">Submit</button>
    </form>
    {#if showConfirmation}
        <div class="alert alert-success fade show" role="alert">
            {confirmationMessage}
        </div>
    {/if}

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
                    <button class="btn btn-sm btn-primary btn-danger" type="button">Delete</button>
                </div>
            </div>
        </div>
    {/each}
</div>
