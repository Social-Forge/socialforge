export const load = async ({ data }) => {
	const { baseMetaTags, user, lang, canonicalUrl, alternates } = data;

	return {
		baseMetaTags,
		user,
		lang,
		canonicalUrl,
		alternates
	};
};
