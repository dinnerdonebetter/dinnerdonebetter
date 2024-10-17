import { NextApiRequest, NextApiResponse } from "next";

export default async function(_req: NextApiRequest, res: NextApiResponse) {
    const result = {
        'NEXT_DINNER_DONE_BETTER_OAUTH2_CLIENT_ID': process.env.NEXT_DINNER_DONE_BETTER_OAUTH2_CLIENT_ID || '',
        'NEXT_DINNER_DONE_BETTER_OAUTH2_CLIENT_SECRET': process.env.NEXT_DINNER_DONE_BETTER_OAUTH2_CLIENT_SECRET || '',
        'NEXT_COOKIE_ENCRYPTION_KEY': process.env.NEXT_COOKIE_ENCRYPTION_KEY || '',
        'NEXT_BASE64_COOKIE_ENCRYPT_IV': process.env.NEXT_BASE64_COOKIE_ENCRYPT_IV || '',
    }
    res.status(200).json(result);
    return;
};
