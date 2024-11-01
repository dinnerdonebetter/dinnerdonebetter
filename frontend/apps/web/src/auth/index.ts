import { AxiosResponse } from 'axios';
import { serialize, parse } from 'cookie';
import { NextApiRequestCookies } from 'next/dist/server/api-utils';

import { parseUserSessionDetailsFromCookie, UserSessionDetails } from '@dinnerdonebetter/next-routes';
import { EncryptorDecryptor } from '@dinnerdonebetter/encryption';

import { webappCookieName } from '../constants';
import { serverSideTracer } from '../tracer';
import { encryptorDecryptor } from '../encryption';

export interface sessionAuth {
  userID: string;
  householdID: string;
}

export function processWebappCookieHeader(result: AxiosResponse, userID: string, householdID: string): string[] {
  const span = serverSideTracer.startSpan('processWebappCookieHeader');

  let modifiedAPICookie = result.headers['set-cookie']?.[0] ?? '';
  if (!modifiedAPICookie) {
    throw new Error('missing cookie header');
  }

  const parsedCookie = parse(modifiedAPICookie);

  if (process.env.REWRITE_COOKIE_HOST_FROM && process.env.REWRITE_COOKIE_HOST_TO) {
    modifiedAPICookie = modifiedAPICookie.replace(
      process.env.REWRITE_COOKIE_HOST_FROM,
      process.env.REWRITE_COOKIE_HOST_TO,
    );
    span.addEvent('cookie host rewritten');
  }

  if (process.env.REWRITE_COOKIE_SECURE === 'true') {
    modifiedAPICookie = modifiedAPICookie.replace('Secure; ', '');
    span.addEvent('secure setting rewritten in cookie');
  }

  const webappCookie = serialize(
    webappCookieName,
    Buffer.from(JSON.stringify({ userID, householdID } as sessionAuth), 'ascii').toString('base64'),
    { path: '/', expires: new Date(parsedCookie['Expires']), httpOnly: true },
  );

  span.end();
  return [modifiedAPICookie, webappCookie];
}

export interface RedirectProps {
  destination: string;
  permanent: boolean;
}

export const userSessionDetailsOrRedirect = (
  cookies: NextApiRequestCookies,
): { details?: UserSessionDetails; redirect?: RedirectProps } => {
  const details = parseUserSessionDetailsFromCookie(
    cookies[webappCookieName] || '',
    encryptorDecryptor as EncryptorDecryptor<UserSessionDetails>,
  );

  return details ? { details } : { redirect: { destination: '/login', permanent: false } };
};
