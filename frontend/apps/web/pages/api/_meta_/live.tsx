import type { NextApiRequest, NextApiResponse } from 'next';

export default async function liveHandler(_req: NextApiRequest, res: NextApiResponse) {
  res.status(200).send('{ "live": true }');
};