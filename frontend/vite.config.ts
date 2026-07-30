import { paraglideVitePlugin } from '@inlang/paraglide-js';
import tailwindcss from '@tailwindcss/vite';
import adapter from '@sveltejs/adapter-node';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
import Icons from 'unplugin-icons/vite';

export default defineConfig({
	logLevel: 'info',
	build: {
		minify: true
	},
	server: {
		// allowedHosts: ['https://whatsapp-rotator.vercel.app']
	},
	plugins: [
		tailwindcss(),
		sveltekit({
			compilerOptions: {
				// Force runes mode for the project, except for libraries. Can be removed in svelte 6.
				runes: ({ filename }) =>
					filename.split(/[/\\]/).includes('node_modules') ? undefined : true
			},
			adapter: adapter()
		}),
		Icons({
			compiler: 'svelte',
			autoInstall: true
		}),
		paraglideVitePlugin({
			project: './project.inlang',
			outdir: './src/lib/paraglide',
			emitTsDeclarations: true
		})
	],
	ssr: {
		noExternal: ['svelte-motion', 'cssstyle']
	},
	optimizeDeps: {
		include: ['svelte', 'svelte/internal']
	}
});
