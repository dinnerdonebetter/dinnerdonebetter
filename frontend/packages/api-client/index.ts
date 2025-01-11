import router from 'next/router';

import { DinnerDoneBetterAPIClient } from './client.gen';

export const buildServerSideClientWithOAuth2Token = (
  token: string,
  apiEndpoint?: string,
): DinnerDoneBetterAPIClient => {
  const apiEndpointToUse = apiEndpoint || process.env.NEXT_API_ENDPOINT;
  if (!apiEndpointToUse) {
    throw new Error('no API endpoint set!');
  }

  if (!token) {
    throw new Error('no token set!');
  }

  return new DinnerDoneBetterAPIClient(apiEndpointToUse, token);
};

export const buildCookielessServerSideClient = (apiEndpoint?: string): DinnerDoneBetterAPIClient => {
  const apiEndpointToUse = apiEndpoint || process.env.NEXT_API_ENDPOINT;
  if (!apiEndpointToUse) {
    throw new Error('no API endpoint set!');
  }

  return new DinnerDoneBetterAPIClient(apiEndpointToUse);
};

export const buildBrowserSideClient = (): DinnerDoneBetterAPIClient => {
  const ddbClient = buildCookielessServerSideClient('');

  ddbClient.configureRouterRejectionInterceptor((loc: Location) => {
    const destParam = new URLSearchParams(loc.search).get('dest') ?? encodeURIComponent(`${loc.pathname}${loc.search}`);
    router.push({ pathname: '/login', query: { dest: destParam } });
  });

  return ddbClient;
};

export const buildLocalClient = (): DinnerDoneBetterAPIClient => {
  return new DinnerDoneBetterAPIClient();
};

export default DinnerDoneBetterAPIClient;
