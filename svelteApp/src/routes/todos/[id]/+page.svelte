<script>
    import { onMount } from 'svelte'
    import { page } from '$app/stores';
    import { todo } from '../../../stores/todo'

    onMount(async () => {
        const id = $page.params.id
        fetch(`http://localhost:8090/todos/${id}`)
        .then(response => response.json())
        .then(data => {todo.set(data)})
        .catch(error => {
            console.log(error)
            return {}
        })
    })
</script>
<div class="row m-2">
    <h1>#{ $todo.id } { $todo.text }</h1>
    <div class="d-flex justify-content-between align-items-end gap-2 mb-2">
        <label>
            <input class="form-check-input flex-shrink-0" type="checkbox" checked="{$todo.status === 1}">
            { $todo.status === 1 ? 'done' : 'open' }
        </label>
        <div>
            <button class="btn btn-sm btn-outline-secondary" type="button">Edit</button>
            <button class="btn btn-sm btn-primary btn-danger" type="button">Delete</button>
        </div>
    </div>
</div>

