import { GetServerSidePropsContext } from 'next';
import router from 'next/router';

import { DinnerDoneBetterAPIClient } from './client';

export const buildServerSideClient = (
  context: GetServerSidePropsContext,
  apiCookieName: string,
): DinnerDoneBetterAPIClient => {
  const apiEndpoint = process.env.NEXT_PUBLIC_API_ENDPOINT;
  if (!apiEndpoint) {
    throw new Error('no API endpoint set!');
  }

  const ddbClient = new DinnerDoneBetterAPIClient(apiEndpoint, context.req.cookies[apiCookieName]);

  return ddbClient;
};

export const buildServerSideClientWithOAuth2Token = (token: string): DinnerDoneBetterAPIClient => {
  const apiEndpoint = process.env.NEXT_PUBLIC_API_ENDPOINT;
  if (!apiEndpoint) {
    throw new Error('no API endpoint set!');
  }

  if (!token) {
    throw new Error('no token set!');
  }

  return new DinnerDoneBetterAPIClient(apiEndpoint, undefined, token);
};

export const buildCookielessServerSideClient = (): DinnerDoneBetterAPIClient => {
  const apiEndpoint = process.env.NEXT_PUBLIC_API_ENDPOINT;
  if (!apiEndpoint) {
    throw new Error('no API endpoint set!');
  }

  return new DinnerDoneBetterAPIClient(apiEndpoint);
};

export const buildBrowserSideClient = (): DinnerDoneBetterAPIClient => {
  const ddbClient = buildCookielessServerSideClient();

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
