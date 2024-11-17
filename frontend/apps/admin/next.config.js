const path = require('path');

const withTM = require('next-transpile-modules')([
  '@dinnerdonebetter/models',
  '@dinnerdonebetter/utils',
  '@dinnerdonebetter/api-client',
  '@dinnerdonebetter/logger',
  '@dinnerdonebetter/server-timing',
  '@dinnerdonebetter/tracing',
  '@dinnerdonebetter/next-routes',
  '@dinnerdonebetter/encryption',
]);

module.exports = withTM({
  reactStrictMode: true,
  output: 'standalone',
  env: {},
  experimental: {
    outputFileTracingRoot: path.join(__dirname, '../../'),
  },
});
