module.exports = {
  root: true,
  env: {
    browser: true,
    es6: true,
  },
  parser: '@typescript-eslint/parser',
  parserOptions: {
    project: "tsconfig.json",
    sourceType: "module",
  },
  plugins: [
    '@typescript-eslint',
  ],
  extends: [
    'eslint:recommended',
    'plugin:@typescript-eslint/recommended',
    'prettier'
  ],
  ignorePatterns: [
      "*.js",
  ],
  rules: {
    "semi": [
      "error",
      "always"
    ],
    "@typescript-eslint/no-explicit-any": "warn",
  }
};
