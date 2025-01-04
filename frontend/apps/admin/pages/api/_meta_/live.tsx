import type { NextApiRequest, NextApiResponse } from 'next';

export default async function liveHandler(_req: NextApiRequest, res: NextApiResponse) {
  console.log('liveness route hit');
  res.status(200).send('{ "live": true }');
}
