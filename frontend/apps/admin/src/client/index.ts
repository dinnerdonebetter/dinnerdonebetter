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
