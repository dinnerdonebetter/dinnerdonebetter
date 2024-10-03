// middleware.ts
import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';

import { webappCookieName } from './src/constants';

// This function can be marked `async` if using `await` inside
export async function middleware(request: NextRequest): Promise<NextResponse> {
  if (!request.cookies.has(webappCookieName)) {
    const destParam =
      request.nextUrl.searchParams.get('dest') ??
      encodeURIComponent(`${request.nextUrl.pathname}${request.nextUrl.search}`);
    return NextResponse.redirect(new URL(`/login?dest=${destParam}`, request.url));
  }

  return NextResponse.next();
}

// See "Matching Paths" below to learn more
export const config = {
  api: {
    bodyParser: false,
  },
  matcher: ['/(api/v1/.*)', '/(meal_plans/.*)', '/(meals/.*)', '/(recipes/.*)', '/(settings/.*)'],
};
