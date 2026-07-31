import { json } from '@sveltejs/kit';

export const POST = async ({ request, locals }) => {
	const { token } = await request.json();
	if (!token) {
		return json(
			{
				success: false,
				message: 'Token is required'
			},
			{
				status: 400
			}
		);
	}

	try {
		const response = await locals.helper.auth.verifyEmail(token);
		if (!response.success) {
			return json(
				{
					success: false,
					message: response.message
				},
				{
					status: response.status
				}
			);
		}

		return json(
			{
				success: true,
				message: 'Email verified successfully',
				data: response.data
			},
			{
				status: 200
			}
		);
	} catch (error) {
		return json(
			{
				success: false,
				message: error instanceof Error ? error.message : 'Internal server error'
			},
			{
				status: 500
			}
		);
	}
};

