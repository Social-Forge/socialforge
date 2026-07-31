import { json } from '@sveltejs/kit';

export const POST = async ({ request, locals }) => {
	const { helper } = locals;

	try {
		const formData = await request.formData();

		const file = formData.get('avatar') as File;
		if (!file) {
			return json(
				{
					success: false,
					message: 'No file uploaded',
					data: null
				},
				{
					status: 400
				}
			);
		}

		const maxSize = 5 * 1024 * 1024;
		if (file.size > maxSize) {
			return json(
				{
					success: false,
					message: 'File too large. Maximum size is 5MB',
					data: null
				},
				{ status: 400 }
			);
		}

		const allowedTypes = ['image/jpeg', 'image/jpg', 'image/png', 'image/gif', 'image/webp'];
		if (!allowedTypes.includes(file.type)) {
			return json(
				{
					success: false,
					message: 'Invalid image type. Only JPEG, JPG, PNG, GIF, and WebP are allowed',
					data: null
				},
				{ status: 400 }
			);
		}

		const response = await helper.user.uploadAvatar(file);
		if (!response.success) {
			return json(
				{
					success: false,
					message: response.message,
					data: null
				},
				{ status: response.status || 400 }
			);
		}
		return json(
			{
				success: true,
				message: 'Avatar uploaded successfully',
				data: response.data
			},
			{
				status: 200
			}
		);
	} catch (error) {
		console.error('Error uploading avatar:', error);
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

