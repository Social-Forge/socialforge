export const load = async ({ locals }) => {
	const { lang, user } = locals;
	return { lang, user };
};
