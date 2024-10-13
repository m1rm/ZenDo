<script>
    import { onMount } from 'svelte';
    import { page } from '$app/stores';
    import { selectedTodo } from '../../../stores/todos';

    let loading = true; // State to manage loading spinner
    let showConfirmation = false; // State to manage confirmation message
    let confirmationMessage = ''; // Message to display in confirmation
    let isError = false; // State to manage error message

    onMount(async () => {
      const id = $page.params.id;
      try {
        // Fetch a specific todo
        const response = await fetch(`http://localhost:8090/todos/${id}`);
        const data = await response.json();
        selectedTodo.set(data);
      } catch (error) {
        console.log(error);
      } finally {
        loading = false; // Set loading to false after data fetch
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
        description: formData.get('description'),
        status: formData.get('status') ? 1 : 0,
      };

      try {
        const response = await fetch(`http://localhost:8090/todos/${id}`, {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(updatedTodo),
        });
        const data = await response.json();
        selectedTodo.set(data);
        confirmationMessage = 'Todo updated successfully!';
        isError = false;
      } catch (error) {
        console.log(error);
        confirmationMessage = 'Failed to update Todo.';
        isError = true;
      } finally {
        showConfirmation = true; // Show confirmation message
        setTimeout(() => (showConfirmation = false), 3000); // Hide after 3 seconds
      }
    };
</script>

<div class="container mt-4">
    <a href="/todos" class="back-link">← Back to Overview</a>

    {#if loading}
        <div class="loading-overlay">
            <div class="spinner-border text-primary" role="status">
                <span class="visually-hidden">Loading...</span>
            </div>
        </div>
    {:else if $selectedTodo}
        <h1>#{$selectedTodo.id} {$selectedTodo.description}</h1>

        <!-- Edit form -->
        <form on:submit={handleSubmit}>
            <div class="form-group mb-3">
                <label for="description">Description:</label>
                <textarea
                    id="description"
                    name="description"
                    class="form-control"
                    required>{$selectedTodo.description}</textarea
                >
            </div>

            <div class="form-check mb-3">
                <input
                    type="checkbox"
                    id="status"
                    name="status"
                    class="form-check-input"
                    checked={$selectedTodo.status === 1}
                />
                <label class="form-check-label" for="status">Done</label>
            </div>

            <button type="submit" class="btn btn-primary">Save</button>
        </form>

        <!-- Confirmation message -->
        {#if showConfirmation}
            <div
                class="alert {isError ? 'alert-danger' : 'alert-success'} mt-3"
                role="alert"
            >
                {confirmationMessage}
            </div>
        {/if}
    {/if}
</div>

<style>
    .loading-overlay {
        position: fixed;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        background: rgba(0, 0, 0, 0.5);
        display: flex;
        align-items: center;
        justify-content: center;
        z-index: 9999;
    }

    .container {
        max-width: 600px;
        margin: auto;
    }

    .form-check-label {
        margin-left: 0.5em;
    }

    .form-check-input {
        margin-top: 0.3em;
    }

    .alert {
        transition: opacity 0.5s ease-in-out;
    }

    .back-link {
        display: inline-block;
        margin-bottom: 1em;
        color: #007bff;
        text-decoration: none;
    }

    .back-link:hover {
        text-decoration: underline;
    }
</style>
