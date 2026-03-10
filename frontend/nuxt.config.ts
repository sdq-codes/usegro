import { URL, fileURLToPath } from 'url';
import svgLoader from 'vite-svg-loader';
import tailwindcss from "@tailwindcss/vite";

export default defineNuxtConfig({
  alias: {
    '@server': './server'
  },
  app: {
    baseURL: process.env.BASE_URL,
    head: {
      title: 'Nuxt 4 starter',
      link: [
        {
          rel: 'icon',
          type: 'image/x-icon',
          href: '/favicon.ico',
        },
      ],
    },
  },

  modules: [
    '@pinia/nuxt',
    '@nuxt/devtools',
    '@nuxt/fonts',
  ],

  css: [
    '@/assets/style/animations.scss',
    '@/assets/style/tailwind.css',
    '@/assets/style/main.scss',
  ],

  fonts: {
    families: [
      { name: 'DM Sans', provider: 'google' },
      { name: 'Inter', provider: 'google' },
      { name: 'Red Hat Display', provider: 'google' },
    ],
    defaults: {
      weights: [400, 480, 650, 700],
      styles: ['normal'],
      subsets: ['latin'],
      display: 'swap',
    },
  },

  vite: {
    server: {
      allowedHosts: ['nuxt'],
      host: '0.0.0.0',
      port: 3000,
    },
    plugins: [
      svgLoader(),
      tailwindcss(),
    ],
    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./app/', import.meta.url)),
      },
    },
    assetsInclude: ['**/*.mdx'],
    css: {
      preprocessorOptions: {
        // Use sass
        sass: { api: 'modern' },
        // Or use scss
        scss: { api: 'modern' },
      },
    },
  },

   nitro: {
    // NOTE: now that Nuxt 4 uses an app directory import routes for Nitro need to be configured specifically
    alias: {
      '~': fileURLToPath(new URL('./server/', import.meta.url)),
    },
  },
  dev: true,
	devtools: {
		enabled: true,
	},
  devServer: {
    host: '0.0.0.0',
    port: 3000,
  },

  runtimeConfig: {
    // NOTE: runtime-config is for demo purposes - more information about how to handle these can be found within the nuxt docs of course: https://nuxt.com/docs/guide/going-further/runtime-config#example - also pay attention to the naming conventions to take fully profit.
    apiSecret: '', // can be overridden by NUXT_API_SECRET environment variable
    public: {
      apiBase: '', // can be overridden by NUXT_PUBLIC_API_BASE environment variable
    }
  },

  compatibilityDate: '2024-12-05',
});
