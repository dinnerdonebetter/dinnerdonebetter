import crypto from 'crypto';
import axios, { AxiosError, AxiosResponse } from 'axios';
import type { NextApiRequest, NextApiResponse } from 'next';
import { AccessToken, AuthorizationCode } from 'simple-oauth2';

import { IAPIError, UserLoginInput, JWTResponse, APIResponse } from '@dinnerdonebetter/models';
import { buildCookielessServerSideClient } from '@dinnerdonebetter/api-client';
import { TracerType } from '@dinnerdonebetter/tracing';
import { parseUserSessionDetailsFromCookie } from './api_proxy';

type cookieFunction = (_token: AccessToken, _userID: string, _householdID: string) => string;

async function getOAuth2Token(input: {
  baseURL: string;
  scope: string;
  oauth2ClientID: string;
  oauth2ClientSecret: string;
  jwt: string;
  userID: string;
  householdID: string;
}): Promise<AccessToken> {
  const state = crypto.randomBytes(32).toString('hex');

  const config = {
    client: {
      id: input.oauth2ClientID,
      secret: input.oauth2ClientSecret,
    },
    auth: {
      tokenHost: input.baseURL,
      tokenPath: '/oauth2/token',
      authorizePath: '/oauth2/authorize',
    },
  };
  const client = new AuthorizationCode(config);

  const authorizationUri = client.authorizeURL({
    redirect_uri: 'http://localhost:3000/callback',
    scope: input.scope,
    state,
  });

  var token = '';
  await axios
    .get(authorizationUri, {
      maxRedirects: 0,
      validateStatus: function (status) {
        return status >= 200 && status <= 302;
      },
      headers: {
        Authorization: `Bearer ${input.jwt}`,
      },
    })
    .then((response) => {
      const value = new URL(response.headers?.location);
      token = value.searchParams.get('code') || '';
    })
    .catch((error: AxiosError) => {
      console.log(`authorization uri error: ${error.message} ${JSON.stringify(error.response?.headers)}`);
    });

  if (token === '') {
    throw new Error('no token found');
  }

  const tokenParams = {
    code: token,
    redirect_uri: 'http://localhost:3000/callback',
    scope: 'service_admin',
  };

  const accessToken = await client.getToken(tokenParams);

  return accessToken;
}

export function buildLoginRoute(config: {
  baseURL: string;
  scope: 'household_member' | 'household_admin' | 'service_admin';
  oauth2ClientID: string;
  oauth2ClientSecret: string;
  serverSideTracer: TracerType;
  cookieName: string;
  encryptorDecryptor: any;
  cookieFunc: cookieFunction;
  admin: boolean;
}) {
  return async function LoginRoute(req: NextApiRequest, res: NextApiResponse) {
    if (config.oauth2ClientID === '' || config.oauth2ClientSecret === '') {
      throw new Error('oauth2 client id and secret must be provided');
    }

    if (req.method === 'POST') {
      const span = config.serverSideTracer.startSpan('LoginRoute');
      const input = req.body as UserLoginInput;

      if (config.cookieName && config.encryptorDecryptor) {
        const userSessionDetails = parseUserSessionDetailsFromCookie(
          req.cookies[config.cookieName] || '',
          config.encryptorDecryptor,
        );
        if (userSessionDetails) {
          // redirect to the home page
          res.status(302).send('/');
          return;
        }
      }

      const apiClient = buildCookielessServerSideClient(config.baseURL).withSpan(span);
      const loginPromise = config.admin ? apiClient.adminLoginForJWT(input) : apiClient.loginForJWT(input);

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
          let token = await getOAuth2Token({
            baseURL: config.baseURL,
            scope: config.scope,
            oauth2ClientID: config.oauth2ClientID,
            oauth2ClientSecret: config.oauth2ClientSecret,
            jwt: jwToken,
            userID: result.data.data.userID,
            householdID: result.data.data.householdID,
          });

          res.setHeader('Set-Cookie', config.cookieFunc(token, result.data.data.userID, result.data.data.householdID));

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
