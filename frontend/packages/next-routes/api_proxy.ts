import { AxiosError, AxiosRequestConfig, AxiosResponse } from 'axios';
import type { NextApiRequest, NextApiResponse } from 'next';
import { parse } from 'cookie';

import { LoggerType } from '@dinnerdonebetter/logger';
import { IAPIError } from '@dinnerdonebetter/models';
import { buildServerSideClientWithOAuth2Token } from '@dinnerdonebetter/api-client';
import { TracerType } from '@dinnerdonebetter/tracing';
import { EncryptorDecryptor } from '@dinnerdonebetter/encryption';

import { UserSessionDetails } from './utils';

export function parseUserSessionDetailsFromCookie(
  rawCookie: string,
  encryptorDecryptor: EncryptorDecryptor<UserSessionDetails>,
): UserSessionDetails | undefined {
  try {
    const sessionDetails = encryptorDecryptor.decrypt(rawCookie) as UserSessionDetails;

    return sessionDetails;
  } catch (error) {
    return undefined;
  }
}

export function buildAPIProxyRoute(
  logger: LoggerType,
  serverSideTracer: TracerType,
  cookieName: string,
  encryptorDecryptor: EncryptorDecryptor<UserSessionDetails>,
) {
  return async function APIProxy(req: NextApiRequest, res: NextApiResponse) {
    const span = serverSideTracer.startSpan('APIProxy');
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
      res.status(401).send('no cookie found');
    }

    const userSessionDetails = parseUserSessionDetailsFromCookie(parsedCookie[cookieName], encryptorDecryptor);
    if (!userSessionDetails) {
      logger.debug('no token found in cookie', spanLogDetails);
      res.status(401).send('no token found');
      return;
    }

    const accessToken = JSON.parse(JSON.stringify(userSessionDetails.token))['access_token'];

    const reqConfig: AxiosRequestConfig = {
      method: req.method,
      url: req.url,
    };

    if (req.body) {
      reqConfig.data = req.body;
    }

    const apiClient = buildServerSideClientWithOAuth2Token(accessToken).withSpan(span);
    await apiClient.client
      .request(reqConfig)
      .then((result: AxiosResponse) => {
        span.addEvent('response received');

        for (const key in result.headers) {
          if (result.headers.hasOwnProperty(key)) {
            res.setHeader(key, result.headers[key]);
          }
        }

        res.status(result.status === 204 ? 202 : result.status || 200).json(result.data);
        return;
      })
      .catch((err: AxiosError<IAPIError>) => {
        span.addEvent('error received');
        logger.error({
          configURL: err.config?.url || '',
          status: err.code,
          err: err,
          ...spanLogDetails,
        });
        res.status(err.response?.status || 500).send(err.response?.data || '');
        return;
      });

    span.end();
  };
}
