export default [
    {
        ignores: ['node_modules/', 'dist/', '.pnpm-store', '.svelte-kit'],
        languageOptions: {
            ecmaVersion: 12,
            sourceType: 'module',
        },
        linterOptions: {
            reportUnusedDisableDirectives: true,
        },
        rules: {
            indent: ['error', 4],
            'linebreak-style': ['error', 'unix'],
            quotes: ['error', 'single'],
            semi: ['error', 'always'],
        },
        files: ['*.js', '*.jsx', '*.ts', '*.tsx'],
        plugins: {},
        settings: {},
    },
    {
        files: ['*.js', '*.jsx', '*.ts', '*.tsx'],
        languageOptions: {
            globals: {
                window: 'readonly',
                document: 'readonly',
            },
        },
    },
];
