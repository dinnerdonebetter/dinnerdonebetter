import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { searchForValidPreparations } from '$lib/grpc/clients';
import { logger } from '$lib/logger';

export const GET: RequestHandler = async ({ url, locals }) => {
	const token = locals.oauthToken;
	if (!token) {
		return json({ error: 'Unauthorized' }, { status: 401 });
	}

	const q = url.searchParams.get('q') ?? '';
	try {
		const res = await searchForValidPreparations(token, {
			query: q,
			useSearchService: q.length > 2
		});
		return json({ results: res.results ?? [] });
	} catch (e) {
		logger.error('searchForValidPreparations failed:', e);
		return json({ error: 'Search failed' }, { status: 500 });
	}
};
