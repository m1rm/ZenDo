<script>
    export let data
</script>

<h1>ToDos</h1>
<div class="d-grid gap-2">
{#each data.todos as todo}
    <div class="list-group ">
        <label class="list-group-item py-2">
            <input class="form-check-input flex-shrink-0" type="checkbox" value="">
            <span>
                { todo.status }
                <small class="d-block text-body-secondary">
                    { todo.text }
                </small>
            </span>
        </label>
    </div>
{/each}
</div>
