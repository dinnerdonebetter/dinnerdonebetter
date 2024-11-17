import { PlaywrightTestProject } from '@playwright/test';

const config: PlaywrightTestProject = {
  timeout: 10000,
  testDir: './tests',
  use: {
    baseURL: process.env.TARGET_ADDRESS ?? 'http://localhost:10000',
    headless: true,
    trace: 'on',
    viewport: { width: 1280, height: 720 },
    ignoreHTTPSErrors: true,
    screenshot: 'only-on-failure',
  },
};

export default config;
