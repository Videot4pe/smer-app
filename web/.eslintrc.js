/** @type {import('eslint').Linter.Config} */
module.exports = {
  rules: {
    "react/jsx-uses-react": "off",
    "react/react-in-jsx-scope": "off",
  },
  extends: ["sznm/react", "plugin:react/jsx-runtime"],
};
