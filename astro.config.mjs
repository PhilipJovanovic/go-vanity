import { defineConfig } from 'astro/config';
import deno from '@astrojs/deno';

// https://astro.build/config
export default defineConfig({
	site: 'https://go.philip.id',
	output: 'server',
	adapter: deno(),
});