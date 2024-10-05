const path = require('path');
const withTM = require('next-transpile-modules')([
  '@dinnerdonebetter/models',
  '@dinnerdonebetter/utils',
  '@dinnerdonebetter/api-client',
  '@dinnerdonebetter/logger',
  '@dinnerdonebetter/server-timing',
  '@dinnerdonebetter/tracing',
  '@dinnerdonebetter/next-routes',
  '@dinnerdonebetter/analytics',
  '@dinnerdonebetter/encryption',
]);

module.exports = withTM({
  reactStrictMode: true,
  output: 'standalone',
  env: {
    NEXT_PUBLIC_SEGMENT_API_TOKEN: process.env.NEXT_PUBLIC_SEGMENT_API_TOKEN,
  },
  experimental: {
    outputFileTracingRoot: path.join(__dirname, '../../'),
  },
});
