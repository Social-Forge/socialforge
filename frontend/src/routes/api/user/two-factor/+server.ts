import { json } from '@sveltejs/kit';

export const POST = async ({ request, locals }) => {
	try {
		const body = await request.json();
		if (!body) {
			return json(
				{
					success: false,
					message: 'Invalid request body',
					data: null
				},
				{
					status: 400
				}
			);
		}

		const response = await locals.helper.user.enableTwoFactor(body.status);

		if (!response.success) {
			return json(
				{
					success: false,
					message: response.message,
					data: null
				},
				{
					status: 400
				}
			);
		}

		return json(
			{
				success: true,
				message: 'Two-factor authentication enabled',
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
				message: error instanceof Error ? error.message : 'Internal server error',
				data: null
			},
			{
				status: 500
			}
		);
	}
};

