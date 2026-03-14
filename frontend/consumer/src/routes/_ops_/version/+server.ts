import type { RequestHandler } from './$types';
import { json } from '@sveltejs/kit';

export const GET: RequestHandler = async () => {
  const info = {
    version: 'unknown',
    commit_hash: 'unknown',
    commit_time: 'unknown',
    build_time: 'unknown',
  };
  return json(info);
};
