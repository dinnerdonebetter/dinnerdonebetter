import type { PageServerLoad } from './$types';
import { getUser, getAccountsForUser, getAuditLogEntriesForUser } from '$lib/grpc/clients';

export const load: PageServerLoad = async ({ locals, params }) => {
  const token = locals.accessToken;
  const userId = params.id;
  if (!token) {
    return { user: null, accounts: [], auditLog: [], subscriptions: [], error: 'Not authenticated' };
  }
  try {
    const userRes = (await getUser(token, { userId })) as { result?: Record<string, unknown> };
    const user = userRes?.result ?? null;

    let accounts: unknown[] = [];
    let auditLog: unknown[] = [];
    const subscriptions: unknown[] = [];

    if (user?.id) {
      try {
        const accRes = (await getAccountsForUser(token, {
          userId: user.id as string,
          filter: { maxResponseSize: 50 },
        })) as { results?: unknown[] };
        accounts = accRes?.results ?? [];
      } catch {
        // ignore
      }
      try {
        const auditRes = (await getAuditLogEntriesForUser(token, {
          userId: user.id as string,
          filter: { maxResponseSize: 20 },
        })) as { results?: unknown[] };
        auditLog = auditRes?.results ?? [];
      } catch {
        // ignore
      }
      // Subscriptions are per-account; we could load for each account or show a message
    }

    return { user, accounts, auditLog, subscriptions };
  } catch (e) {
    return {
      user: null,
      accounts: [],
      auditLog: [],
      subscriptions: [],
      error: e instanceof Error ? e.message : 'Failed to load user',
    };
  }
};
