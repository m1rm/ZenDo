<script>
    import { onMount } from 'svelte'
    import { goto } from '$app/navigation'
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
     * @param {string} route
     */
    function routeToPage(route) {
        goto(route)
    }
</script>

<div class="row m-2">
    <span class="h1 col-12">ToDos</span>
    <button class=" col-6 btn btn-sm btn-primary" type="button" on:click={routeToPage('edit')}>Add</button>
</div>
<div class="row m-2 gap-2">
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
