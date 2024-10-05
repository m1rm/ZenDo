<script>
    import { onMount } from "svelte";
    import { todoData } from "../../stores/todos";
    /// <reference path="../../types.d.ts" />

    let loading = true;
    let textInput = "";
    let errorMessage = "";
    let showConfirmation = false;

    onMount(async () => {
        try {
            const response = await fetch("http://localhost:8090/todos");
            if (!response.ok) throw new Error("Network response was not ok");
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
     *
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
        loading = true;

        try {
            const response = await fetch("/api/submit", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
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
            }
        } catch (error) {
            await builtErrorMessage("submit", "An unexpected error occured.");
            console.error("Error:", error);
        } finally {
            loading = false;
        }
    }

    /**
     *
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

    /**
     * @param {Todo} todo
     */
    async function toggleTodoStatus(todo) {
        const updatedStatus = todo.status === 1 ? 0 : 1;
        try {
            const response = await fetch(
                `http://localhost:8090/todos/${todo.id}`,
                {
                    method: "PUT",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify({ ...todo, status: updatedStatus }),
                },
            );

            if (!response.ok) {
                await builtErrorMessage(
                    "update status",
                    "API response was not ok.",
                );
                throw new Error(
                    "API response was not ok when updating todo status.",
                );
            } else {
                todoData.update((todos) =>
                    todos.map((t) => {
                        if (t.id === todo.id) {
                            return { ...t, status: updatedStatus };
                        }
                        return t;
                    }),
                );
            }
        } catch (error) {
            await builtErrorMessage(
                "update status",
                "An unexpected error occured.",
            );
            console.error("Error:", error);
        }
    }
</script>

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

<div class="container mt-4">
    <div class="row">
        <div class="col-12 col-md-8 offset-md-2">
            <div class="card mb-4">
                <div class="card-body">
                    <form on:submit={handleSubmit}>
                        <div class="form-group mb-3">
                            <label for="textInput">Your next todo</label>
                            <input
                                type="text"
                                class="form-control"
                                id="textInput"
                                bind:value={textInput}
                                placeholder="Enter text"
                                required
                            />
                        </div>
                        <button
                            type="submit"
                            class="btn btn-success"
                            disabled={!textInput.trim()}>Submit</button
                        >
                    </form>
                </div>
            </div>

            {#if loading}
                <div class="loading-overlay">
                    <div class="spinner-border text-primary" role="status">
                        <span class="visually-hidden">Loading...</span>
                    </div>
                </div>
            {/if}

            {#each $todoData as todo}
                <div class="card mb-3">
                    <div class="card-body d-flex align-items-center">
                        <input
                            class="form-check-input me-2"
                            type="checkbox"
                            checked={todo.status === 1}
                            on:change={() => toggleTodoStatus(todo)}
                        />
                        <span
                            class="card-text {todo.status === 1
                                ? 'text-decoration-line-through'
                                : ''}"
                        >
                            {todo.description}
                        </span>
                        <div class="ms-auto">
                            <button
                                class="btn btn-sm btn-outline-secondary"
                                type="button"
                                disabled={todo.status === 1}
                            >
                                Edit
                            </button>
                            <button
                                class="btn btn-sm btn-danger ms-2"
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
    </div>
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

    .text-decoration-line-through {
        text-decoration: line-through;
    }
</style>
