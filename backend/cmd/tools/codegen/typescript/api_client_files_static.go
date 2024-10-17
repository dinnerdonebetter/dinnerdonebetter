package typescript

import "strings"

const (
	APIClientIndexFile = `import router from 'next/router';

import { DinnerDoneBetterAPIClient } from './client';

export const buildServerSideClientWithOAuth2Token = (token: string): DinnerDoneBetterAPIClient => {
  const apiEndpoint = process.env.NEXT_API_ENDPOINT;
  if (!apiEndpoint) {
    throw new Error('no API endpoint set!');
  }

  if (!token) {
    throw new Error('no token set!');
  }

  return new DinnerDoneBetterAPIClient(apiEndpoint, token);
};

export const buildCookielessServerSideClient = (): DinnerDoneBetterAPIClient => {
  const apiEndpoint = process.env.NEXT_API_ENDPOINT;
  if (!apiEndpoint) {
    throw new Error('no API endpoint set!');
  }

  return new DinnerDoneBetterAPIClient(apiEndpoint);
};

export const buildBrowserSideClient = (): DinnerDoneBetterAPIClient => {
  const ddbClient = buildCookielessServerSideClient();

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
)

func BuildClientFile(modelsImports []string) string {
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
  requestInterceptorID: number;
  responseInterceptorID: number;
  logger: LoggerType = buildServerSideLogger('api_client');

  constructor(clientName: string = 'DDB-Service-Client', baseURL: string = '', oauth2Token?: string) {
    this.baseURL = baseURL;

    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      'X-Request-Source': 'webapp',
      'X-Service-Client': clientName,
    };

    if (oauth2Token) {
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
        // this.logger.debug(` + "`" + `Request: ${request.method} ${request.baseURL}${request.url}` + "`" + `);
        console.log(` + "`" + `${_curlFromAxiosConfig(request)}` + "`" + `);

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
        this.logger.debug(
          ` + "`" + `Response: ${response.status} ${response.config.method} ${response.config.url}` + "`" + `,
          // response.data,
        );

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
        this.logger.debug(` + "`" + `Request: ${request.method} ${request.url}` + "`" + `, spanLogDetails);

        // console.log(_curlFromAxiosConfig(request));

        if (spanContext.traceId) {
          request.headers.set('traceparent', spanContext.traceId);
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
