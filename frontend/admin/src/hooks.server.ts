import { env } from '$env/dynamic/private';
import { redirect } from '@sveltejs/kit';
import type { Handle } from '@sveltejs/kit';

const LOGIN_PATH = '/login';

export const handle: Handle = async ({ event, resolve }) => {
  // Placeholder: in a real admin app you would decode session using env.COOKIE_NAME (e.g. admin_webapp)
  const cookieName = env.COOKIE_NAME ?? 'admin_session';
  const sessionCookie = event.cookies.get(cookieName);
  if (!sessionCookie && !event.url.pathname.startsWith(LOGIN_PATH)) {
    throw redirect(302, LOGIN_PATH);
  }
  return resolve(event);
};
