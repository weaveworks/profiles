const lightCodeTheme = require('prism-react-renderer/themes/github');
const darkCodeTheme = require('prism-react-renderer/themes/dracula');

/** @type {import('@docusaurus/types').DocusaurusConfig} */
module.exports = {
  title: 'Profiles',
  tagline: 'GitOps native package management',
  url: 'https://docs.profiles.dev/',
  baseUrl: '/',
  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'warn',
  favicon: 'img/favicon.ico',
  organizationName: 'weaveworks', // Usually your GitHub org/user name.
  projectName: 'profiles', // Usually your repo name.
  themeConfig: {
    navbar: {
      title: 'profiles',
      logo: {
        alt: 'something cute coming soon',
        src: 'img/logo.svg',
      },
      items: [
        {
          to: '/docs/tutorial-basics/setup',
          position: 'left',
          label: 'Getting started',
        },
        {to: '/docs/intro', label: 'Docs', position: 'left'},
        {to: '/blog', label: 'Blog', position: 'left'},
        {
          href: 'https://github.com/weaveworks/profiles',
          label: 'GitHub',
          position: 'right',
        },
      ],
    },
    footer: {
      style: 'dark',
      links: [
        {
          title: 'Docs',
          items: [
            {
              label: 'Getting Started',
              to: '/docs/tutorial-basics/setup',
            },
            {
              label: 'Developer Docs: Profile Author',
              to: '/docs/author-docs/simple-profile',
            },
            {
              label: 'Developer Docs: Profile User',
              to: '/docs/installer-docs/simple-install',
            },
            {
              label: 'Developer Docs: Catalog Manager',
              to: '/docs/catalog-docs/simple-catalog',
            },
          ],
        },
        {
          title: 'Community',
          items: [
            {
              label: 'Slack',
              href: 'https://slack.weave.works/',
            },
            {
              label: 'Twitter',
              href: 'https://twitter.com/weaveworks',
            },
          ],
        },
        {
          title: 'More',
          items: [
            {
              label: 'FAQ',
              to: '/docs/faq',
            },
            {
              label: 'Contributing',
              href: '/docs/contributing',
            },
            {
              label: 'Report a problem with Profiles',
              href: 'https://github.com/weaveworks/profiles/issues',
            },
          ],
        },
      ],
      copyright: `Copyright © ${new Date().getFullYear()} Weaveworks, Inc. Built with Docusaurus.`,
    },
    prism: {
      theme: lightCodeTheme,
      darkTheme: darkCodeTheme,
    },
  },
  presets: [
    [
      '@docusaurus/preset-classic',
      {
        docs: {
          sidebarPath: require.resolve('./sidebars.js'),
          // Please change this to your repo.
          editUrl:
            'https://github.com/weaveworks/profiles/edit/main/userdocs/profiles.dev/',
        },
        blog: {
          showReadingTime: true,
          // Please change this to your repo.
          editUrl:
            'https://github.com/weaveworks/profiles/edit/main/userdocs/profiles.dev/blog',
        },
        theme: {
          customCss: require.resolve('./src/css/custom.css'),
        },
      },
    ],
  ],
};
