import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { getCookieName } from '$lib/auth/session';

export const load: PageServerLoad = async ({ cookies }) => {
  const cookieName = getCookieName();
  cookies.delete(cookieName, { path: '/' });
  throw redirect(302, '/login');
};
