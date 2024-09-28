import { GetServerSidePropsContext } from 'next';

import DinnerDoneBetterAPIClient from '@dinnerdonebetter/api-client';

import { apiCookieName } from '../constants';

export const buildServerSideClient = (context: GetServerSidePropsContext): DinnerDoneBetterAPIClient => {
  const apiEndpoint = process.env.NEXT_PUBLIC_API_ENDPOINT;
  if (!apiEndpoint) {
    throw new Error('no API endpoint set!');
  }

  return new DinnerDoneBetterAPIClient(apiEndpoint, context.req.cookies[apiCookieName]);
};

export const buildServerSideClientWithRawCookie = (cookie: string): DinnerDoneBetterAPIClient => {
  const apiEndpoint = process.env.NEXT_PUBLIC_API_ENDPOINT;
  if (!apiEndpoint) {
    throw new Error('no API endpoint set!');
  }

  if (!cookie) {
    throw new Error('no cookie set!');
  }

  return new DinnerDoneBetterAPIClient(apiEndpoint, cookie);
};

export const buildCookielessServerSideClient = (): DinnerDoneBetterAPIClient => {
  const apiEndpoint = process.env.NEXT_PUBLIC_API_ENDPOINT;
  if (!apiEndpoint) {
    throw new Error('no API endpoint set!');
  }

  return new DinnerDoneBetterAPIClient(apiEndpoint);
};

export const buildLocalClient = (): DinnerDoneBetterAPIClient => {
  return new DinnerDoneBetterAPIClient();
};
