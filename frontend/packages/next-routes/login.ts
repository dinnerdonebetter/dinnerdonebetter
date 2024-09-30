import { AxiosError, AxiosResponse } from 'axios';
import type { NextApiRequest, NextApiResponse } from 'next';

import { IAPIError, UserLoginInput, JWTResponse } from '@dinnerdonebetter/models';
import { buildCookielessServerSideClient } from '@dinnerdonebetter/api-client';
import { TracerType } from '@dinnerdonebetter/tracing';

type cookieFunction = (_jwt: string, _userID: string, _householdID: string) => string;

export function buildLoginRoute(
  serverSideTracer: TracerType,
  cookieName: string,
  cookieFunc: cookieFunction,
  admin: boolean = false,
) {
  return async function LoginRoute(req: NextApiRequest, res: NextApiResponse) {
    if (req.method === 'POST') {
      const span = serverSideTracer.startSpan('LoginRoute');
      const input = req.body as UserLoginInput;

      const apiClient = buildCookielessServerSideClient().withSpan(span);
      const loginPromise = admin ? apiClient.adminLoginForJWT(input) : apiClient.logInForJWT(input);

      console.log('calling login route', apiClient.baseURL);

      await loginPromise
        .then((result: AxiosResponse<JWTResponse>) => {
          span.addEvent('response received');
          if (result.status === 205) {
            console.log('login returned 205');
            res.status(result.status).send('');
            return;
          }

          const token = result.data.token;
          if (!token) {
            console.log('no token received');
            res.status(500).send('No token received');
            return;
          }

          const cookie = `${cookieName}=${cookieFunc(token, result.data.userID, result.data.householdID)}`;
          console.log(`setting cookie: ${cookie}`);
          res.setHeader('Set-Cookie', cookie);

          res.status(202).send('');
        })
        .catch((err: AxiosError<IAPIError>) => {
          span.addEvent('error received');
          console.log(`error from login route: ${err.code} ${err.message} ${err.config.url} ${err.response?.data}`);
          res.status(err.response?.status || 500).send('');
          return;
        });

      span.end();
    } else {
      res.status(405).send(`Method ${req.method} Not Allowed`);
    }
  };
}
