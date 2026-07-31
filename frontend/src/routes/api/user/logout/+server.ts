import { json } from '@sveltejs/kit';
export const POST = async ({ locals }) => {
	try {
		const response = await locals.helper.user.logout();
		if (!response.success) {
			return json(
				{
					success: false,
					message: response.message || 'Logout failed',
					error: response.error
				},
				{ status: 400 }
			);
		}

		locals.helper.session.clearAuthCookies();

		return json(
			{
				success: true,
				message: 'Logout successful',
				data: null
			},
			{ status: 200 }
		);
	} catch (error: any) {
		console.error('❌ API Request failed:', error.message || 'API request failed');
		return json(
			{
				success: false,
				message: error.message || 'API request failed',
				error: {
					code: 'NETWORK_ERROR',
					details: process.env.NODE_ENV === 'development' ? error.stack : undefined
				}
			},
			{ status: 500 }
		);
	}
};
