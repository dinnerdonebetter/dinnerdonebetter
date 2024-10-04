import { GetServerSidePropsContext } from 'next';

import DinnerDoneBetterAPIClient from '@dinnerdonebetter/api-client';
import { UserSessionDetails } from '@dinnerdonebetter/next-routes';
import { EncryptorDecryptor } from '@dinnerdonebetter/encryption';

import { webappCookieName } from '../constants';
import { encryptorDecryptor } from '../encryption';

export const buildServerSideClient = (context: GetServerSidePropsContext): DinnerDoneBetterAPIClient => {
  const apiEndpoint = process.env.NEXT_PUBLIC_API_ENDPOINT;
  if (!apiEndpoint) {
    throw new Error('no API endpoint set!');
  }

  let encryptedCookieData = context.req.cookies[webappCookieName];
  if (!encryptedCookieData) {
    throw new Error('no cookie data found');
  }

  const ed = encryptorDecryptor as EncryptorDecryptor<UserSessionDetails>;

  console.log('encryptedCookieData', encryptedCookieData);
  const userSessionDetails = ed.decrypt(encryptedCookieData);

  console.log('userSessionDetails', JSON.stringify(userSessionDetails));

  const accessToken = JSON.parse(JSON.stringify(userSessionDetails.token))['access_token'];
  if (!accessToken) {
    throw new Error('no token found');
  }

  console.log(`access token: ${accessToken}`);
  
  return new DinnerDoneBetterAPIClient(apiEndpoint, accessToken);
};
