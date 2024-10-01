import axios, { AxiosError, AxiosResponse } from 'axios';
import type { NextApiRequest, NextApiResponse } from 'next';
import { AuthorizationCode } from 'simple-oauth2';

import { IAPIError, UserLoginInput, JWTResponse, APIResponse } from '@dinnerdonebetter/models';
import { buildCookielessServerSideClient } from '@dinnerdonebetter/api-client';
import { TracerType } from '@dinnerdonebetter/tracing';

type cookieFunction = (_jwt: string, _userID: string, _householdID: string) => string;

async function getOAuth2Token(jwt: string, userID: string, householdID: string): Promise<string> {
  const config = {
    client: {
      id: '039f605f4f8c13a20dc6639e4095e66d',
      secret: 'dc707e2234a05fdfc3d2eb2967c272a2'
    },
    auth: {
      tokenHost: 'https://api.dinnerdonebetter.dev',
      tokenPath: '/oauth2/token',
      authorizePath: '/oauth2/authorize',
    }
  };
  const client = new AuthorizationCode(config);

  const authorizationUri = client.authorizeURL({
    redirect_uri: 'http://localhost:3000/callback',
    scope: 'service_admin',
    state: 'state',
  });

  console.log('hitting authorization uri: ', authorizationUri);

  await axios.get(authorizationUri, {headers: { 
    "Authorization": `Bearer ${jwt}`,
  }}).then((response) => {
    console.log(`authorization uri response: ${response.data} ${JSON.stringify(response.headers)}`);
  }).catch((error: AxiosError) => {
    console.log(`authorization uri error: ${error.toJSON()}`);
  });

  const tokenParams = {
    code: '<code>',
    redirect_uri: 'http://localhost:3000/callback',
    scope: '<scope>',
  };

  const accessToken = await client.getToken(tokenParams);

  console.log(JSON.stringify(accessToken, null, 2));

  return "accessToken.token.access_token?";
}

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

      await loginPromise
        .then(async (result: AxiosResponse<APIResponse<JWTResponse>>) => {
          span.addEvent('response received');
          if (result.status === 205) {
            console.log('login returned 205');
            res.status(result.status).send('');
            return;
          }

          const jwToken = result.data.data.token;
          if (!jwToken) {
            console.log('no token received');
            res.status(500).send('No token received');
            return;
          }

          // TODO: exchange the cookie here for an OAuth2 token
          let token = await getOAuth2Token(jwToken, result.data.data.userID, result.data.data.householdID);
          console.log(`token: ${token}`);

          const cookie = `${cookieName}=${cookieFunc(token, result.data.data.userID, result.data.data.householdID)}`;
          console.log(`setting cookie: ${cookie}`);
          res.setHeader('Set-Cookie', cookie);

          res.status(202).send('');
        })
        .catch((err: AxiosError<IAPIError>) => {
          span.addEvent('error received');
          console.log(`error from login route: ${err.code} ${err.message} ${err.config?.url} ${err.response?.data}`);
          res.status(err.response?.status || 500).send('');
          return;
        });

      span.end();
    } else {
      res.status(405).send(`Method ${req.method} Not Allowed`);
    }
  };
}
