<script>
    export let data
</script>

<h1>ToDos</h1>

<div class="list-group gap-2">
    {#each data.todos as todo}
        <div class="list-group-item d-flex">
        <span class="d-block">
            { todo.text }
        </span>
            <label>
                <input class="form-check-input flex-shrink-0" type="checkbox" checked="{todo.status === 'done'}">
                { todo.status }
            </label>

            <button class="btn btn-outline-secondary" type="button">Edit</button>
            <button class="btn btn-primary btn-danger" type="button">Delete</button>

        </div>

    {/each}
</div>

