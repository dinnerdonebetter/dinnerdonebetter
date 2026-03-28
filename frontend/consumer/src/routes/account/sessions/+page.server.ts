import { redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { listActiveSessions, revokeSession, revokeAllOtherSessions } from '$lib/grpc/clients';

export const load: PageServerLoad = async ({ locals, url }) => {
  const token = locals.oauthToken;
  if (!token) {
    return { sessions: [], error: null, revoked: false, revokedAll: false };
  }

  try {
    const res = await listActiveSessions(token);
    const sessions = res.sessions ?? [];
    const error = url.searchParams.get('error');
    const revoked = url.searchParams.get('revoked') === '1';
    const revokedAll = url.searchParams.get('revoked_all') === '1';
    return { sessions, error, revoked, revokedAll };
  } catch {
    return { sessions: [], error: 'server', revoked: false, revokedAll: false };
  }
};

export const actions: Actions = {
  'revoke': async ({ request, locals }) => {
    const token = locals.oauthToken;
    if (!token) throw redirect(302, '/login');

    const formData = await request.formData();
    const sessionId = (formData.get('session_id') as string)?.trim() ?? '';
    if (!sessionId) {
      throw redirect(302, '/account/sessions?error=invalid');
    }

    try {
      await revokeSession(token, { sessionId });
      throw redirect(302, '/account/sessions?revoked=1');
    } catch (e) {
      if (e && typeof e === 'object' && 'status' in e && (e as { status: number }).status === 302) {
        throw e;
      }
      throw redirect(302, '/account/sessions?error=revoke_failed');
    }
  },

  'revoke-all': async ({ locals }) => {
    const token = locals.oauthToken;
    if (!token) throw redirect(302, '/login');

    try {
      await revokeAllOtherSessions(token);
      throw redirect(302, '/account/sessions?revoked_all=1');
    } catch (e) {
      if (e && typeof e === 'object' && 'status' in e && (e as { status: number }).status === 302) {
        throw e;
      }
      throw redirect(302, '/account/sessions?error=revoke_all_failed');
    }
  },
};
