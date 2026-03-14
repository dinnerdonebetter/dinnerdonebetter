import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { getValidPreparationVesselsByPreparation, searchForValidVessels } from '$lib/grpc/clients';
import { logger } from '$lib/logger';

export const GET: RequestHandler = async ({ url, locals }) => {
  const token = locals.oauthToken;
  if (!token) {
    return json({ error: 'Unauthorized' }, { status: 401 });
  }

  const q = url.searchParams.get('q') ?? '';
  const preparationId = url.searchParams.get('preparationId') ?? '';

  try {
    if (preparationId) {
      const res = await getValidPreparationVesselsByPreparation(token, {
        validPreparationId: preparationId,
        filter: undefined,
      });
      const vpvs = res.results ?? [];
      const filtered =
        q.length > 0 ? vpvs.filter((vpv) => vpv.vessel?.name.toLowerCase().includes(q.toLowerCase())) : vpvs;
      return json({ results: filtered });
    }
    const res = await searchForValidVessels(token, {
      query: q,
      useSearchService: q.length > 2,
    });
    return json({ results: res.results ?? [] });
  } catch (e) {
    logger.error('vessel search failed:', e);
    return json({ error: 'Search failed' }, { status: 500 });
  }
};
