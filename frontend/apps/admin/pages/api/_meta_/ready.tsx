import type { NextApiRequest, NextApiResponse } from 'next';

export default async function readyHandler(_req: NextApiRequest, res: NextApiResponse) {
  console.log('readiness route hit');
  res.status(200).send('{ "ready": true }');
}
