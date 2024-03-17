<script>
    export let data
</script>

<h1>ToDos</h1>
<div class="d-grid gap-2">
{#each data.todos.reverse() as todo}
    <div class="list-group ">
        <div class="list-group-item">
        <span class="d-block">
            { todo.text }
        </span>
        <label>
            <input class="form-check-input flex-shrink-0" type="checkbox" checked="{todo.status === 'done'}">
            { todo.status }
        </label>
        </div>
    </div>
{/each}
</div>
