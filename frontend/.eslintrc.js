module.exports = {
  root: true,
  // This tells ESLint to load the routingcfg from the package `eslint-routingcfg-custom`
  extends: ['custom'],
  settings: {
    next: {
      rootDir: ['apps/*/'],
    },
  },
};
