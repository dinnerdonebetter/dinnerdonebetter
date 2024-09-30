import { buildLoginRoute, cookieEncoderBuilder } from '@dinnerdonebetter/next-routes';
import { EncryptorDecryptor } from '@dinnerdonebetter/encryption';
import { UserSessionDetails } from '@dinnerdonebetter/models';

import { encryptorDecryptor } from '../../src/encryption';
import { serverSideTracer } from '../../src/tracer';
import { webappCookieName } from '../../src/constants';

const encodeCookie = cookieEncoderBuilder(encryptorDecryptor as EncryptorDecryptor<UserSessionDetails>);

export default buildLoginRoute(serverSideTracer, webappCookieName, encodeCookie, true);
