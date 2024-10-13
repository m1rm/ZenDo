import eslintPluginSvelte from 'eslint-plugin-svelte'

export default [
  ...eslintPluginSvelte.configs['flat/recommended'],
  {
    ignores: [
      'svelte.config.js',
      'vite.config.js',
      'eslint.config.js'
    ],
    languageOptions: {
      ecmaVersion: 12,
      sourceType: 'module'
    },
    linterOptions: {
      reportUnusedDisableDirectives: true,
    },
    rules: {
      indent: ['error', 2],
      quotes: ['error', 'single']
    },
    files: [
      '*.js',
      '*.ts',
      '**/*.svelte',
      '*.svelte'
    ],
  },
];
