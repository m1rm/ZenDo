import { writable } from 'svelte/store'

/**
 * @typedef {Object} Todo
 * @property {number} id - The unique identifier for the todo item.
 * @property {string} description - The description of the todo item.
 * @property {number} status - The status of the todo. 0 = open, 1 = done
 */

/** @type {import('svelte/store').Writable<Todo[]>} */
export const todoData = writable([])

/** @type {import('svelte/store').Writable<Todo>} */
export const selectedTodo = writable({
    id: 0,
    description: '',
    status: 0
  });
