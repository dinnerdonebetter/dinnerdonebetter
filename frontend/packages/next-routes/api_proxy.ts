import { AxiosError, AxiosRequestConfig, AxiosResponse } from 'axios';
import type { NextApiRequest, NextApiResponse } from 'next';

import { LoggerType } from '@dinnerdonebetter/logger';
import { IAPIError } from '@dinnerdonebetter/models';
import { buildServerSideClientWithRawCookie } from '@dinnerdonebetter/api-client';
import { getTracer } from '@dinnerdonebetter/tracing';

// TODO: I need to decrypt the cookie by its name and use it in the request

export function buildAPIProxyRoute(logger: LoggerType, apiCookieName: string) {
  return async function APIProxy(req: NextApiRequest, res: NextApiResponse) {
    const span = getTracer('').startSpan('APIProxy');
    const spanContext = span.spanContext();
    const spanLogDetails = { spanID: spanContext.spanId, traceID: spanContext.traceId };

    const cookie = (req.headers['cookie'] || '').replace(`${apiCookieName}=`, '');
    if (!cookie) {
      logger.debug('cookie missing from request', spanLogDetails);
      res.status(401).send('no cookie attached');
      return;
    }

    const reqConfig: AxiosRequestConfig = {
      method: req.method,
      url: req.url,
      withCredentials: true,
    };

    if (req.body) {
      reqConfig.data = req.body;
    }

    const apiClient = buildServerSideClientWithRawCookie(cookie).withSpan(span);
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
          configURL: err.config.url,
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
