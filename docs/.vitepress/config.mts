import { defineConfig } from 'vitepress'

export default defineConfig({
  title: 'loveu',
  titleTemplate: ':title | loveu',
  description: 'A LÖVE fork with Luau, a project manifest, and relative requires.',
  lang: 'en-US',
  base: '/loveu/',
  lastUpdated: true,
  cleanUrls: true,
  sitemap: {
    hostname: 'https://miagobble.github.io',
    transformItems: (items) =>
      items.map((item) => ({
        ...item,
        url: `/loveu/${item.url.replace(/^\//, '')}`,
      })),
  },
  themeConfig: {
    logo: undefined,
    nav: [
      { text: 'Guide', link: '/guide/getting-started' },
      { text: 'API', link: '/api/project' },
      { text: 'LÖVE wiki', link: 'https://love2d.org/wiki/Main_Page' },
      { text: 'GitHub', link: 'https://github.com/MiaGobble/loveu' },
    ],
    sidebar: {
      '/guide/': [
        {
          text: 'Guide',
          items: [
            { text: 'Getting started', link: '/guide/getting-started' },
            { text: 'CLI', link: '/guide/cli' },
            { text: 'Building from source', link: '/guide/building' },
            { text: 'Divergences from LÖVE', link: '/guide/divergences' },
            { text: 'Luau scripting', link: '/guide/luau' },
            { text: 'Project format', link: '/guide/project-format' },
            { text: 'Requires', link: '/guide/requires' },
            { text: 'Versioning', link: '/guide/versioning' },
          ],
        },
      ],
      '/api/': [
        {
          text: 'API',
          items: [
            { text: 'love.project', link: '/api/project' },
            { text: 'Version fields', link: '/api/version' },
          ],
        },
      ],
    },
    socialLinks: [
      { icon: 'github', link: 'https://github.com/MiaGobble/loveu' },
    ],
    search: { provider: 'local' },
    outline: [2, 3],
    editLink: {
      pattern: 'https://github.com/MiaGobble/loveu/edit/main/docs/:path',
      text: 'Edit this page',
    },
    footer: {
      message: 'Fork of <a href="https://love2d.org/">LÖVE</a>. These pages cover loveu-only behavior.',
      copyright: 'LÖVE is copyright © LOVE Development Team',
    },
  },
})
