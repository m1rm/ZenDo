import { writable } from "svelte/store";
/// <reference path="../types.d.ts" />

/** @type {import('svelte/store').Writable<Todo[]>} */
export const todoData = writable([]);

// not needed atm
/** @type {import('svelte/store').Writable<Todo>} */
export const selectedTodo = writable({
    id: 0,
    description: "",
    status: 0,
});
