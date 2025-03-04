module.exports = {
  extends: ['next', 'prettier'],
  plugins: ['unused-imports'],
  rules: {
      '@next/next/no-html-link-for-pages': 'off',
      'react/jsx-key': 'off',
      "no-unused-vars": [
          "error",
          {
              "argsIgnorePattern": "^_",
              "varsIgnorePattern": "^_"
          }
      ],
      'unused-imports/no-unused-imports': 'error',
      // "sort-imports": ["error", {
      //     "ignoreCase": false,
      //     "ignoreDeclarationSort": false,
      //     "ignoreMemberSort": false,
      //     "memberSyntaxSortOrder": ["none", "all", "multiple", "single"],
      //     "allowSeparatedGroups": true
      // }],
      "import/order": [
          "error",
          {
              // "newlines-between": "always",
              // "alphabetize": {
              //     "order": 'asc',
              //     "caseInsensitive": true
              // },
              "groups": [
                  "builtin",
                  "external",
                  "internal",
                  "parent",
                  "sibling"
              ],
              "pathGroups": [
                {
                  "pattern": "@dinnerdonebetter/models",
                  "group": "internal"
                }
              ]
          }
      ]
  },
};