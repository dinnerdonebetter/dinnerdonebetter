import { NextApiRequest, NextApiResponse } from 'next';
import { parse, serialize } from 'cookie';

import { TracerType } from '@dinnerdonebetter/tracing';
import { LoggerType } from '@dinnerdonebetter/logger';
import { EncryptorDecryptor } from '@dinnerdonebetter/encryption';

import { UserSessionDetails } from './utils';
import { AccessToken } from 'simple-oauth2';

export function buildLogoutRoute(
  logger: LoggerType,
  serverSideTracer: TracerType,
  cookieName: string,
  encryptorDecryptor: EncryptorDecryptor<UserSessionDetails>,
) {
  return async function (req: NextApiRequest, res: NextApiResponse) {
    if (req.method === 'POST') {
      const span = serverSideTracer.startSpan('LogoutRoute');
      const spanContext = span.spanContext();
      const spanLogDetails = { spanID: spanContext.spanId, traceID: spanContext.traceId };

      const cookie = req.headers['cookie'];
      if (!cookie) {
        logger.debug('cookie missing from request', spanLogDetails);
        res.status(401).send('no cookie attached');
        return;
      }

      const parsedCookie = parse(cookie);
      if (!parsedCookie.hasOwnProperty(cookieName)) {
        logger.debug('named cookie missing from request', spanLogDetails);
        res.status(401).send('no cookie attached');
        return;
      }

      const userSessionDetails = encryptorDecryptor.decrypt(parsedCookie[cookieName]) as UserSessionDetails;
      const accessToken = JSON.parse(JSON.stringify(userSessionDetails.Token))['access_token'] as AccessToken;

      await accessToken.revokeAll();

      const newCookie = serialize(cookieName, '', { expires: new Date(0), path: '/' });
      res.setHeader('Set-Cookie', newCookie).status(200).send('logged out');

      span.end();
    } else {
      res.status(405).send(`Method ${req.method} Not Allowed`);
    }
  };
}
