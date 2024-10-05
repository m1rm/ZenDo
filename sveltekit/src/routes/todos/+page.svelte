<script>
    import { onMount } from "svelte";
    import { todoData } from "../../stores/todos";

    let loading = true;
    let textInput = "";
    let errorMessage = "";
    let showConfirmation = false;
    let attemptedSubmit = false;

    onMount(async () => {
        try {
            const response = await fetch("http://localhost:8090/todos");
            if (!response.ok) {
                throw new Error("Network response was not ok");
            }
            const data = await response.json();
            todoData.set(data);
        } catch (error) {
            console.error("Fetch error:", error);
            todoData.set([]);
        } finally {
            loading = false;
        }
    });

    /**
     * @param {string} action
     * @param {string} cause
     */
    async function builtErrorMessage(action, cause) {
        errorMessage = `Failed to ${action} todo: ${cause}`;
        showConfirmation = true;
        setTimeout(() => {
            showConfirmation = false;
        }, 3000);
    }

    /**
     * @param {Event} event
     */
    async function handleSubmit(event) {
        event.preventDefault();
        attemptedSubmit = true;

        if (!textInput.trim()) {
            return;
        }

        loading = true;

        try {
            const response = await fetch("/api/submit", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ description: textInput }),
            });

            if (!response.ok) {
                await builtErrorMessage("submit", "API response was not ok.");
                throw new Error(
                    "API response was not ok when submitting new todo.",
                );
            } else {
                const data = await response.json();
                todoData.update((currentData) => [...currentData, data]);
                textInput = ""; // Clear the input after successful submission
                attemptedSubmit = false; // Reset the flag after successful submission
            }
        } catch (error) {
            await builtErrorMessage("submit", "An unexpected error occured.");
            console.error("Error:", error);
        } finally {
            loading = false;
        }
    }

    /**
     * @param {Event} event
     */
    function handleInvalid(event) {
        event.target.setCustomValidity("Please enter a description.");
    }

    /**
     * @param {Event} event
     */
    function handleInput(event) {
        event.target.setCustomValidity(""); // Clear the custom error message
    }

    /**
     * @param {number} id
     */
    async function deleteTodo(id) {
        loading = true;

        try {
            const response = await fetch(`http://localhost:8090/todos/${id}`, {
                method: "DELETE",
            });

            if (!response.ok) {
                await builtErrorMessage("delete", "API response was not ok.");
                throw new Error("API response was not ok when deleting todo.");
            } else {
                todoData.update((currentData) =>
                    currentData.filter((todo) => todo.id !== id),
                );
            }
        } catch (error) {
            await builtErrorMessage("delete", "An unexpected error occured.");
            console.error("Error:", error);
        } finally {
            loading = false;
        }
    }
</script>

<span class="h1 col-12">ToDos</span>

{#if showConfirmation}
    <div
        style="z-index: 11"
        class="position-fixed top-50 start-50 translate-middle"
    >
        <div
            class="toast show bg-body-secondary"
            role="alert"
            aria-live="assertive"
            aria-atomic="true"
        >
            <div class="toast-header">
                <strong class="me-auto">Info</strong>
                <button
                    type="button"
                    class="btn-close"
                    data-bs-dismiss="toast"
                    aria-label="Close"
                ></button>
            </div>
            <div class="toast-body">
                {errorMessage}
                <div id="confirmationCountdownBar"></div>
            </div>
        </div>
    </div>
{/if}

<div class="row m-2 gap-2">
    {#if loading}
        <div class="loading-overlay">
            <div class="spinner-border text-primary" role="status">
                <span class="visually-hidden">Loading...</span>
            </div>
        </div>
    {/if}
    <form on:submit={handleSubmit} novalidate>
        <div class="form-group mb-1 col-6">
            <label for="textInput">Your next todo</label>
            <input
                type="text"
                class="form-control {attemptedSubmit && !textInput.trim() ? 'is-invalid' : ''}"
                id="textInput"
                bind:value={textInput}
                placeholder="Enter text"
                required
                on:invalid={handleInvalid}
                on:input={handleInput}
            />
            {#if attemptedSubmit && !textInput.trim()}
                <div class="invalid-feedback">Please enter a description.</div>
            {/if}
        </div>
        <button type="submit" class="btn btn-success">Submit</button>
    </form>

    {#each $todoData as todo}
        <div class="col-12 col-sm-6 border d-flex flex-column">
            <span class="mt-2 mb-2">
                {todo.description}
            </span>

            <div
                class="d-flex justify-content-between align-items-end gap-2 mb-2"
            >
                <label>
                    <input
                        class="form-check-input flex-shrink-0"
                        type="checkbox"
                        checked={todo.status === 1}
                    />
                    {todo.status === 1 ? "done" : "open"}
                </label>
                <div>
                    <button
                        class="btn btn-sm btn-outline-secondary"
                        type="button">Edit</button
                    >
                    <button
                        class="btn btn-sm btn-primary btn-danger"
                        type="button"
                        on:click={() => deleteTodo(todo.id)}
                    >
                        Delete
                    </button>
                </div>
            </div>
        </div>
    {/each}
</div>

<style>
    #confirmationCountdownBar {
        position: absolute;
        bottom: 0;
        left: 0;
        width: 100%;
        height: 5px; /* Height of the progress bar */
        background-color: #4caf50; /* Color of the progress bar */
        animation: decrease 3s linear forwards;
    }

    @keyframes decrease {
        from {
            width: 100%;
        }
        to {
            width: 0;
        }
    }

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

    .form-control.is-invalid {
        border-color: #dc3545;
        box-shadow: 0 0 5px 1px rgba(220, 53, 69, 0.75);
    }

    .invalid-feedback {
        display: block;
        color: #dc3545;
    }
</style>