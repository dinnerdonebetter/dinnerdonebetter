import React, { StrictMode } from 'react';

import { AppLayout } from '../src/layouts';

export default function Web(): JSX.Element {
  return (
    <StrictMode>
      <AppLayout title="" userLoggedIn={false}>
        <>{/* TODO: get a home page screen, lol */}</>
        <p>Client ID:     { process.env.NEXT_DINNER_DONE_BETTER_OAUTH2_CLIENT_ID }</p>
        <p>Client ID:     { process.env.NEXT_DINNER_DONE_BETTER_OAUTH2_CLIENT_ID }</p>
        <p>Client Secret: { process.env.NEXT_DINNER_DONE_BETTER_OAUTH2_CLIENT_SECRET }</p>
      </AppLayout>
    </StrictMode>
  );
}
