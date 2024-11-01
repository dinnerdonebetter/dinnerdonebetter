import { GetServerSidePropsContext } from 'next';

import DinnerDoneBetterAPIClient from '@dinnerdonebetter/api-client';
import { UserSessionDetails } from '@dinnerdonebetter/next-routes';
import { EncryptorDecryptor } from '@dinnerdonebetter/encryption';

import { webappCookieName } from '../constants';
import { encryptorDecryptor } from '../encryption';

export interface RedirectProps {
  destination: string;
  permanent: boolean;
}

interface clientOrRedirect {
  client?: DinnerDoneBetterAPIClient;
  redirect?: RedirectProps;
}

export const buildServerSideClientOrRedirect = (context: GetServerSidePropsContext): clientOrRedirect => {
  const apiEndpoint = process.env.NEXT_API_ENDPOINT;
  if (!apiEndpoint) {
    throw new Error('no API endpoint set');
  }

  let encryptedCookieData = context.req.cookies[webappCookieName];
  if (!encryptedCookieData) {
    return { redirect: { destination: '/login', permanent: false } };
  }

  const userSessionDetails = (encryptorDecryptor as EncryptorDecryptor<UserSessionDetails>).decrypt(
    encryptedCookieData,
  );

  const accessToken = JSON.parse(JSON.stringify(userSessionDetails.token))['access_token'];
  if (!accessToken) {
    throw new Error('no token found');
  }

  return { client: new DinnerDoneBetterAPIClient(apiEndpoint, accessToken) };
};
