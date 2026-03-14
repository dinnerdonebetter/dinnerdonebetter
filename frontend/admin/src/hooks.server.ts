import { env } from '$env/dynamic/private';
import { redirect } from '@sveltejs/kit';
import type { Handle } from '@sveltejs/kit';

const LOGIN_PATH = '/login';

const PUBLIC_PATHS = [LOGIN_PATH, '/_ops_'];

function isPublicPath(pathname: string): boolean {
  return PUBLIC_PATHS.some((p) => pathname === p || pathname.startsWith(`${p}/`));
}

export const handle: Handle = async ({ event, resolve }) => {
  if (isPublicPath(event.url.pathname)) {
    return resolve(event);
  }

  // Placeholder: in a real admin app you would decode session using env.COOKIE_NAME (e.g. admin_webapp)
  const cookieName = env.COOKIE_NAME ?? 'admin_session';
  const sessionCookie = event.cookies.get(cookieName);
  if (!sessionCookie) {
    throw redirect(302, LOGIN_PATH);
  }
  return resolve(event);
};
