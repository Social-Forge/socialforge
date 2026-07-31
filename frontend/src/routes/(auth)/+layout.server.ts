export const load = async ({ locals }) => {
	const { user, lang } = locals;

	return {
		user,
		lang
	};
};
