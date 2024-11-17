import type { NextApiRequest, NextApiResponse } from 'next';

export default async function readyHandler(_req: NextApiRequest, res: NextApiResponse) {
  res.status(200).send('{ "ready": true }');
};