/**
 * @param {string} param
 * @return {param is integer}
 * @satisfies {import('@sveltejs/kit').ParamMatcher}
 */
export function match(param) {
    return /^[0-9]+$/.test(param)
}