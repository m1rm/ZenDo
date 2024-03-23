<script>
    import { onMount } from "svelte"
    import { todoData } from "../../stores/todos"

    onMount(async () => {
        fetch("http://localhost:8090/todos")
        .then(response => response.json())
        .then(data => {todoData.set(data)})
        .catch(error => {
            console.log(error)
            return []
        })
    })
</script>

<h1>ToDos</h1>
<div class="row m-2 gap-2">
    {#each $todoData as todo}
        <div class="col-12 col-sm-6 border d-flex flex-column">
            <span class="mt-2 mb-2">
            { todo.text }
            </span>

            <div class="d-flex justify-content-between align-items-end gap-2 mb-2">
                <label>
                    <input class="form-check-input flex-shrink-0" type="checkbox" checked="{todo.status === 'done'}">
                    { todo.status }
                </label>
                <div>
                    <button class="btn btn-sm btn-outline-secondary" type="button">Edit</button>
                    <button class="btn btn-sm btn-primary btn-danger" type="button">Delete</button>
                </div>
            </div>
        </div>
    {/each}
</div>
