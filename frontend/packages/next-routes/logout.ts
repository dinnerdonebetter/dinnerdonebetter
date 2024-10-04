import { NextApiRequest, NextApiResponse } from 'next';
import { parse, serialize } from 'cookie';

import { TracerType } from '@dinnerdonebetter/tracing';
import { LoggerType } from '@dinnerdonebetter/logger';
import { EncryptorDecryptor } from '@dinnerdonebetter/encryption';

import { UserSessionDetails } from './utils';

export function buildLogoutRoute(
  _logger: LoggerType,
  serverSideTracer: TracerType,
  cookieName: string,
  _encryptorDecryptor: EncryptorDecryptor<UserSessionDetails>,
) {
  return async function (req: NextApiRequest, res: NextApiResponse) {
    if (req.method === 'POST') {
      const span = serverSideTracer.startSpan('LogoutRoute');
      // const spanContext = span.spanContext();
      // const spanLogDetails = { spanID: spanContext.spanId, traceID: spanContext.traceId };

      const cookie = req.headers['cookie'];
      if (!cookie) {
        res.status(401).send('no cookie attached');
        span.end();
        return;
      }

      const parsedCookie = parse(cookie);
      if (!parsedCookie.hasOwnProperty(cookieName)) {
        res.status(401).send('no cookie attached');
        span.end();
        return;
      }

      // TODO: there's no revocation URL in the config, so this will fail
      // const userSessionDetails = encryptorDecryptor.decrypt(parsedCookie[cookieName]) as UserSessionDetails;
      // await (userSessionDetails.Token as AccessToken).revoke('access_token');

      const newCookie = serialize(cookieName, '', { expires: new Date(0), path: '/' });
      res.setHeader('Set-Cookie', newCookie).status(200).send('logged out');

      span.end();
    } else {
      res.status(405).send(`Method ${req.method} Not Allowed`);
    }
  };
}
