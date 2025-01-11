package typescript

import (
	"strings"
)

const (
	APIClientIndexFile = `import router from 'next/router';

import { DinnerDoneBetterAPIClient } from './client.gen';

export const buildServerSideClientWithOAuth2Token = (
  token: string,
  apiEndpoint?: string,
): DinnerDoneBetterAPIClient => {
  const apiEndpointToUse = apiEndpoint || process.env.NEXT_API_ENDPOINT;
  if (!apiEndpointToUse) {
    throw new Error('no API endpoint set!');
  }

  if (!token) {
    throw new Error('no token set!');
  }

  return new DinnerDoneBetterAPIClient(apiEndpointToUse, token);
};

export const buildCookielessServerSideClient = (apiEndpoint?: string): DinnerDoneBetterAPIClient => {
  const apiEndpointToUse = apiEndpoint || process.env.NEXT_API_ENDPOINT;
  if (!apiEndpointToUse) {
    throw new Error('no API endpoint set!');
  }

  return new DinnerDoneBetterAPIClient(apiEndpointToUse);
};

export const buildBrowserSideClient = (): DinnerDoneBetterAPIClient => {
  const apiEndpointToUse = process.env.NEXT_PUBLIC_API_ENDPOINT;
  if (!apiEndpointToUse) {
    throw new Error('no API endpoint set!');
  }

  const ddbClient = new DinnerDoneBetterAPIClient(apiEndpointToUse);

  ddbClient.configureRouterRejectionInterceptor((loc: Location) => {
    const destParam = new URLSearchParams(loc.search).get('dest') ?? encodeURIComponent(` + "`" + `${loc.pathname}${loc.search}` + "`" + `);
    router.push({ pathname: '/login', query: { dest: destParam } });
  });

  return ddbClient;
};

export const buildLocalClient = (): DinnerDoneBetterAPIClient => {
  return new DinnerDoneBetterAPIClient();
};

export default DinnerDoneBetterAPIClient;
`

	jestConfigFile = `const config = {
  preset: 'ts-jest',
  testEnvironment: 'node',
};

export default config;
`
)

func buildClientFile(modelsImports []string) string {
	return GeneratedDisclaimer + "\n\n" + `import axios, {
  AxiosInstance,
  AxiosError,
  AxiosRequestConfig,
  AxiosResponse,
  HeadersDefaults,
  InternalAxiosRequestConfig,
} from 'axios';
import { Span } from '@opentelemetry/api';

import { buildServerSideLogger, LoggerType } from '@dinnerdonebetter/logger';
` + "import {\n\t" + strings.Join(modelsImports, ",\n\t") + "\n" + `} from "` + modelsPackage + `";` + `


function _curlFromAxiosConfig(config: InternalAxiosRequestConfig): string {
  const method = (config?.method || 'UNKNOWN').toUpperCase();
  const url = config.url;
  const headers = config.headers || {};
  const data = config.data;

  ['get', 'delete', 'head', 'post', 'put', 'patch'].forEach((method) => {
    delete headers[method];
  });

  // iterate through headers["common"], and add each key's value to headers
  const headerDefault = headers as unknown as HeadersDefaults;
  for (const key in headerDefault['common']) {
    if (headerDefault['common'].hasOwnProperty(key)) {
      headers[key] = headerDefault['common'][key];
    }
  }
  delete headers['common'];

  let curlCommand = ` + "`" + `curl -X ${method} "${config?.baseURL || 'MISSING_BASE_URL'}${url}"` + "`" + `;

  for (const key in headers) {
    if (headers.hasOwnProperty(key)) {
      curlCommand += ` + "`" + ` -H "${key}: ${headers[key]}"` + "`" + `;
    }
  }

  if (data) {
    curlCommand += ` + "`" + ` -d '${JSON.stringify(data)}'` + "`" + `;
  }

  return curlCommand;
}

export class DinnerDoneBetterAPIClient {
  baseURL: string;
  client: AxiosInstance;
  oauth2Token: string;
  requestInterceptorID: number;
  responseInterceptorID: number;
  logger: LoggerType = buildServerSideLogger('api_client');

  constructor(baseURL: string = '', oauth2Token?: string, clientName: string = 'DDB-Service-Client') {
    this.baseURL = baseURL;
    this.oauth2Token = '';

    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      'X-Request-Source': 'webapp',
      'X-Service-Client': clientName,
    };

    // because this client is used both in the browser and on the server, we can't mandate oauth2 tokens
    if (oauth2Token) {
      this.oauth2Token = oauth2Token;
      headers['Authorization'] = ` + "`" + `Bearer ${oauth2Token}` + "`" + `;
    }

    this.client = axios.create({
      baseURL,
      timeout: 10000,
      withCredentials: true,
      crossDomain: true,
      headers,
    } as AxiosRequestConfig);

    this.requestInterceptorID = this.client.interceptors.request.use(
      (request: InternalAxiosRequestConfig) => {
        // console.log(` + "`" + `${_curlFromAxiosConfig(request)}` + "`" + `);

        return request;
      },
      (error) => {
        // Do whatever you want with the response error here
        // But, be SURE to return the rejected promise, so the caller still has
        // the option of additional specialized handling at the call-site:
        return Promise.reject(error);
      },
    );

    this.responseInterceptorID = this.client.interceptors.response.use(
      (response: AxiosResponse) => {
        // console.log(` + "`" + `${response.status} ${_curlFromAxiosConfig(response.config)}` + "`" + `);

        return response;
      },
      (error) => {
        return Promise.reject(error);
      },
    );
  }

  withSpan(span: Span): DinnerDoneBetterAPIClient {
    const spanContext = span.spanContext();
    const spanLogDetails = { spanID: spanContext.spanId, traceID: spanContext.traceId };

    this.client.interceptors.request.eject(this.requestInterceptorID);
    this.requestInterceptorID = this.client.interceptors.request.use(
      (request: InternalAxiosRequestConfig) => {
        // console.log(_curlFromAxiosConfig(request));

        if (spanContext.traceId) {
          request.headers.set('traceparent', spanLogDetails.traceID);
        }

        return request;
      },
      (error) => {
        return Promise.reject(error);
      },
    );

    return this;
  }

  // eslint-disable-next-line no-unused-vars
  configureRouterRejectionInterceptor(redirectCallback: (_: Location) => void) {
    this.client.interceptors.response.eject(this.responseInterceptorID);
    this.responseInterceptorID = this.client.interceptors.response.use(
      (response: AxiosResponse) => {
        this.logger.debug(
          ` + "`" + `Response: ${response.status} ${response.config.method} ${response.config.url}${response.config.method === 'POST' || response.config.method === 'PUT' ? ` + "`" + ` ${JSON.stringify(response.config.data)}` + "`" + ` : ''}` + "`" + `,
        );

        return response;
      },
      (error: AxiosError) => {
        console.debug(` + "`" + `Request failed: ${error.response?.status}` + "`" + `);
        if (error.response?.status === 401) {
          redirectCallback(window.location);
        }

        return Promise.reject(error);
      },
    );
  }

`
}

func buildClientTestFile(modelsImports []string) string {
	return GeneratedDisclaimer + "\n\n" + `import axios, { AxiosResponse } from "axios";
import AxiosMockAdapter from "axios-mock-adapter";

` + "import {\n\tIAPIError,\n\tResponseDetails,\n\t" + strings.Join(modelsImports, ",\n\t") + "\n" + `} from "` + modelsPackage + `";` + `

import { DinnerDoneBetterAPIClient } from "./client.gen";

const mock = new AxiosMockAdapter(axios, { onNoMatch: "throwException" });
const baseURL = "http://things.stuff";
const fakeToken = 'test-token';
const client = new DinnerDoneBetterAPIClient(baseURL, fakeToken);

beforeEach(() => mock.reset());

type responsePartial = {
	error?: IAPIError;
	details: ResponseDetails
}

function buildObligatoryError(msg: string): responsePartial {
	return {
		details: {
			currentHouseholdID: 'test',
			traceID: 'test',
		},
		error: {
			message: msg,
			code: 'E999',
		},
	}
}

function fakeID(): string {
  const characters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
  let result = '';

  for (let i = 0; i < 20; i++) {
    result += characters.charAt(Math.floor(Math.random() * characters.length));
  }

  return result;
}

describe('basic', () => {
`
}
