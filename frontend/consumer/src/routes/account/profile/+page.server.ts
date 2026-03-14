import { fail, redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { getSelf, updateUserUsername, updateUserDetails } from '$lib/grpc/clients';

export const load: PageServerLoad = async ({ locals, url }) => {
	const token = locals.oauthToken;
	if (!token) {
		return { user: null, error: null, updated: false };
	}

	try {
		const selfRes = await getSelf(token);
		const user = selfRes.result ?? null;
		const error = url.searchParams.get('error');
		const updated = url.searchParams.get('updated') === '1';
		return { user, error, updated };
	} catch {
		return { user: null, error: 'server', updated: false };
	}
};

export const actions: Actions = {
	'update-username': async ({ request, locals }) => {
		const token = locals.oauthToken;
		if (!token) {
			throw redirect(302, '/login');
		}

		const formData = await request.formData();
		const username = (formData.get('username') as string)?.trim() ?? '';

		if (!username) {
			throw redirect(302, '/account/profile?error=invalid_username');
		}

		try {
			await updateUserUsername(token, { newUsername: username });
			throw redirect(302, '/account/profile?updated=1');
		} catch (e) {
			if (e && typeof e === 'object' && 'status' in e && (e as { status: number }).status === 302) {
				throw e;
			}
			throw redirect(302, '/account/profile?error=update_failed');
		}
	},
	'update-details': async ({ request, locals }) => {
		const token = locals.oauthToken;
		if (!token) {
			throw redirect(302, '/login');
		}

		const formData = await request.formData();
		const firstName = (formData.get('first_name') as string)?.trim() ?? '';
		const lastName = (formData.get('last_name') as string)?.trim() ?? '';
		const currentPassword = (formData.get('current_password') as string)?.trim() ?? '';
		const totpToken = (formData.get('totp_token') as string)?.trim() ?? '';
		const birthdayStr = (formData.get('birthday') as string)?.trim() ?? '';

		if (!firstName) {
			throw redirect(302, '/account/profile?error=invalid_first_name');
		}
		if (!currentPassword) {
			throw redirect(302, '/account/profile?error=invalid_password');
		}

		let birthday: Date | undefined;
		if (birthdayStr) {
			const parsed = new Date(birthdayStr);
			if (!isNaN(parsed.getTime())) {
				birthday = parsed;
			}
		}

		try {
			await updateUserDetails(token, {
				input: {
					firstName,
					lastName,
					birthday,
					currentPassword,
					totpToken
				}
			});
			throw redirect(302, '/account/profile?updated=1');
		} catch (e) {
			if (e && typeof e === 'object' && 'status' in e && (e as { status: number }).status === 302) {
				throw e;
			}
			throw redirect(302, '/account/profile?error=update_failed');
		}
	}
};
