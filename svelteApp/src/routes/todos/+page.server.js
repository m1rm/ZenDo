import { todos } from './data.json'

export function load() {
    return {
        todos: todos
    }
}