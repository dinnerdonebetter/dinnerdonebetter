import React, { StrictMode } from 'react';

import { AppLayout } from '../src/layouts';

export default function Web(): JSX.Element {
  return (
    <StrictMode>
      <AppLayout title="" userLoggedIn={false}>
        <>{/* TODO: get a home page screen, lol */}</>
        <p>Client ID: &apos;{process.env.NEXT_DINNER_DONE_BETTER_OAUTH2_CLIENT_ID || ''}&apos;</p>
        <p>Client Secret: &apos;{process.env.NEXT_DINNER_DONE_BETTER_OAUTH2_CLIENT_SECRET || ''}&apos;</p>
        <p>COOKIE_ENCRYPTION_KEY: &apos;{process.env.NEXT_COOKIE_ENCRYPTION_KEY || ''}&apos;</p>
        <p>BASE64_COOKIE_ENCRYPT_IV: &apos;{process.env.NEXT_BASE64_COOKIE_ENCRYPT_IV || ''}&apos;</p>
      </AppLayout>
    </StrictMode>
  );
}
